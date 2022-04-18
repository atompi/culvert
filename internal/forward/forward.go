package forward

import (
	"fmt"
	"io"
	"net"
	"os"

	"gitee.com/autom-studio/culvert/internal/config"
	"golang.org/x/crypto/ssh"
)

func Tunnel(sshClientConfig *ssh.ClientConfig, sshServerConfig *config.SSHServerConfig, culvertConfig config.CulvertConfig) {
	for _, tunnel := range culvertConfig.Tunnels {
		config.Wg.Add(1)
		go forward(tunnel, sshClientConfig, sshServerConfig)
	}
	config.Wg.Wait()
}

func forward(tunnel config.TunnelConfig, sshClientConfig *ssh.ClientConfig, sshServerConfig *config.SSHServerConfig) {
	defer config.Wg.Done()
	localListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", tunnel.Local.Bind, tunnel.Local.Port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Create local listener failed:", err)
		os.Exit(1)
	}

	for {
		sshClientConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshServerConfig.Host, sshServerConfig.Port), sshClientConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Start ssh client connection failed:", err)
		}
		remoteConn, err := sshClientConn.Dial("tcp", fmt.Sprintf("%s:%d", tunnel.Remote.Host, tunnel.Remote.Port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connect to remote failed:", err)
			os.Exit(1)
		}

		localConn, err := localListener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Accept Local listener failed:", err)
			os.Exit(1)
		}

		config.Wg.Add(1)
		go func() {
			defer config.Wg.Done()
			_, err = io.Copy(remoteConn, localConn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "io.Copy remote to local failed:", err)
			}
		}()
		config.Wg.Add(1)
		go func() {
			defer config.Wg.Done()
			_, err = io.Copy(localConn, remoteConn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "io.Copy local to remote failed:", err)
			}
		}()
	}
}
