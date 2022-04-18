# Culvert

A tool to quickly open ssh forwarding tunnels.

## Build

```
git clone https://gitee.com/autom-studio/culvert.git
go build -o examples/culvert
```

## Start

### modify config yaml

```
cd examples
cp culvert.yaml.example culvert.yaml
```

Content of `culvert.yaml`:

```
---
server:
  host: 192.168.15.128    # SSH server host
  port: 22    # SSH server port
  username: atompi    # SSH server login user
  password: "123456"    # SSH server login password, if empty use private key
  keyFile: "./id_rsa"    # SSH server login private key, if empty use password, if both empty print error and exit

tunnels:    # tunnels list
  - name: devdocs_128    # tunnel name
    remote:
      host: 192.168.15.128    # Forwarded remote host
      port: 9292    # Forwarded remote port
    local:
      bind: 0.0.0.0    # bind local host
      port: 19292    # bind local port
  - name: ssh_128
    remote:
      host: 192.168.15.128
      port: 22
    local:
      bind: 192.168.15.128
      port: 2222
```

### start

+ Option 1: Start frontend

```
./culvert
```

+ Option 2: Start with Systemd

```
mkdir -p /app/culvert
cp examples/culvert /app/culvert/
cp examples/culvert.yaml /app/culvert/
cp examples/culvert.service /lib/systemd/system/
systemctl daemon-reload
systemctl start culvert
```
