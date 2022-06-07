package tunnel

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

func readFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func (t TunnelConfig) BindTunnel(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var once sync.Once // Only print errors once per session
		var sshAuth []ssh.AuthMethod
		func() {
			// Connect to the tunnel host via SSH.
			// Generate sshAuth
			if t.Host.Password != "" {
				sshAuth = []ssh.AuthMethod{
					ssh.Password(t.Host.Password),
				}
			} else if t.Host.KeyFile != "" {
				if t.Host.KeyPassword != "" {
					key, err := readFile(t.Host.KeyFile)
					if err != nil {
						zap.L().Sugar().Errorf("read private key failed: %v", err)
					}
					signer, err := ssh.ParsePrivateKey(key)
					if err != nil {
						zap.L().Sugar().Errorf("parse private key failed: %v", err)
					}
					sshAuth = []ssh.AuthMethod{
						ssh.PublicKeys(signer),
					}
				} else {
					key, err := readFile(t.Host.KeyFile)
					if err != nil {
						zap.L().Sugar().Errorf("read private key failed: %v", err)
					}
					signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(t.Host.KeyPassword))
					if err != nil {
						zap.L().Sugar().Errorf("parse private key failed: %v", err)
					}
					sshAuth = []ssh.AuthMethod{
						ssh.PublicKeys(signer),
					}
				}
			}
			// Generate hostKeys
			_, _, pubKey, _, _, err := ssh.ParseKnownHosts([]byte(t.Host.KnownHost))
			if err != nil {
				zap.L().Sugar().Errorf("parse known hosts failed: %v", err)
			}
			hostKeys := ssh.FixedHostKey(pubKey)
			// Create ssh dial
			hostAddr := fmt.Sprintf("%s:%d", t.Host.IP, t.Host.Port)
			dial, err := ssh.Dial("tcp", hostAddr, &ssh.ClientConfig{
				User:            t.Host.Username,
				Auth:            sshAuth,
				HostKeyCallback: hostKeys,
				Timeout:         5 * time.Second,
			})
			if err != nil {
				once.Do(func() { zap.L().Sugar().Errorf("create ssh dial error: %v", err) })
				return
			}

			wg.Add(1)
			go t.keepAliveMonitor(&once, wg, dial)
			defer dial.Close()

			// Attempt to bind to the inbound socket.
			var l net.Listener
			switch t.Mode {
			case "L":
				l, err = net.Listen("tcp", fmt.Sprintf("%s:%d", t.Local.Bind, t.Local.Port))
			case "R":
				l, err = dial.Listen("tcp", fmt.Sprintf("%s:%d", t.Remote.Bind, t.Remote.Port))
			}
			if err != nil {
				once.Do(func() { zap.L().Sugar().Errorf("listener bind error: %v", err) })
				return
			}

			// The socket is binded. Make sure we close it eventually.
			bindCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			go func() {
				dial.Wait()
				cancel()
			}()
			go func() {
				<-bindCtx.Done()
				once.Do(func() {}) // Suppress future errors
				l.Close()
			}()

			zap.L().Sugar().Infof("tunnel %s listener bind", t.Name)
			defer zap.L().Sugar().Infof("tunnel %s listener collapsed", t.Name)

			// Accept all incoming connections.
			for {
				conn, err := l.Accept()
				if err != nil {
					once.Do(func() { zap.L().Sugar().Errorf("tunnel %s listener accept error: %v", t.Name, err) })
					return
				}
				wg.Add(1)
				go t.dialTunnel(bindCtx, wg, dial, conn)
			}
		}()

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(t.RetryInterval) * time.Second):
			zap.L().Sugar().Warnf("%s retrying...", t.Name)
		}
	}
}

func (t TunnelConfig) dialTunnel(ctx context.Context, wg *sync.WaitGroup, client *ssh.Client, conn net.Conn) {
	defer wg.Done()

	// The inbound connection is established. Make sure we close it eventually.
	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-connCtx.Done()
		conn.Close()
	}()

	// Establish the outbound connection.
	var dialConn net.Conn
	var err error
	switch t.Mode {
	case "L":
		dialConn, err = client.Dial("tcp", fmt.Sprintf("%s:%d", t.Remote.Bind, t.Remote.Port))
	case "R":
		dialConn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", t.Local.Bind, t.Local.Port))
	}
	if err != nil {
		zap.L().Sugar().Errorf("tunnel %s dial error: %v", t.Name, err)
		return
	}

	go func() {
		<-connCtx.Done()
		dialConn.Close()
	}()

	zap.L().Sugar().Infof("tunnel %s connection established", t.Name)
	defer zap.L().Sugar().Infof("tunnel %s connection closed", t.Name)

	// Copy bytes from one connection to the other until one side closes.
	var once sync.Once
	var dialWg sync.WaitGroup
	dialWg.Add(2)
	go func() {
		defer dialWg.Done()
		defer cancel()
		if _, err := io.Copy(conn, dialConn); err != nil {
			once.Do(func() { zap.L().Sugar().Errorf("tunnel %s send data error: %v", t.Name, err) })
		}
		once.Do(func() {}) // Suppress future errors
	}()
	go func() {
		defer dialWg.Done()
		defer cancel()
		if _, err := io.Copy(dialConn, conn); err != nil {
			once.Do(func() { zap.L().Sugar().Errorf("tunnel %s receive data error: %v", t.Name, err) })
		}
		once.Do(func() {}) // Suppress future errors
	}()
	dialWg.Wait()
}

// keepAliveMonitor periodically sends messages to invoke a response.
// If the server does not respond after some period of time,
// assume that the underlying net.Conn abruptly died.
func (t TunnelConfig) keepAliveMonitor(once *sync.Once, wg *sync.WaitGroup, client *ssh.Client) {
	defer wg.Done()
	if t.Keepalive.Interval == 0 || t.Keepalive.CountMax == 0 {
		return
	}

	// Detect when the SSH connection is closed.
	wait := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		wait <- client.Wait()
	}()

	// Repeatedly check if the remote server is still alive.
	var aliveCount int32
	ticker := time.NewTicker(time.Duration(t.Keepalive.Interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case err := <-wait:
			if err != nil && err != io.EOF {
				once.Do(func() { zap.L().Sugar().Errorf("tunnel %s connection error: %v", t.Name, err) })
			}
			return
		case <-ticker.C:
			if n := atomic.AddInt32(&aliveCount, 1); n > int32(t.Keepalive.CountMax) {
				once.Do(func() { zap.L().Sugar().Warnf("tunnel %s ssh keepalive termination", t.Name) })
				client.Close()
				return
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
			if err == nil {
				atomic.StoreInt32(&aliveCount, 0)
			}
		}()
	}
}
