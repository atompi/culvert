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
tunnels:
  - name: devdocs_128    # tunnel name
    mode: L    # tunnel mode L/R like ssh -L / ssh -R
    host:
      ip: 192.168.15.128    # tunnel host ip
      port: 22    tunnel host ssh port
      username: atompi    # tunnel host ssh login user
      password: "123456"    # tunnel host ssh login password
      keyFile: "./id_rsa"    # tunnel host ssh login private key
      keyPassword: ""    # tunnel host ssh login private key password
      knownHost: "192.168.15.128 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBOUDn8tF9i1XwSnYYKnoyR9z4g+pgdMR16vFFVH1UpskxgpAjgjBubdqTmIs1JQ8OJyWBomqandNM2WtIgQqAPc="    # known_hosts for tunnel host, generate by command: ssh-keyscan -t ecdsa -p 22 192.168.15.128
    keepalive:
      interval: 30    # send keepalive package interval
      countMax: 2    # max send count
    remote:
      bind: 192.168.15.128    # remote bind ip
      port: 9292    # remote bind port
    local:
      bind: 0.0.0.0    # local bind ip
      port: 19292    # local bind port
    retryInterval: 5    retry connection interval
  - name: ssh_128
    mode: R
    host:
      ip: 192.168.15.128
      port: 22
      username: atompi
      password: "123456"
      keyFile: "./id_rsa"
      keyPassword: ""
      knownHost: "192.168.15.128 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBOUDn8tF9i1XwSnYYKnoyR9z4g+pgdMR16vFFVH1UpskxgpAjgjBubdqTmIs1JQ8OJyWBomqandNM2WtIgQqAPc="
    keepalive:
      interval: 30
      countMax: 2
    remote:
      bind: 0.0.0.0    # only support 0.0.0.0 when create tunnel via "R" mode
      port: 2222
    local:
      bind: 192.168.15.128
      port: 22
    retryInterval: 5

log:
  path: "./culvert.log"
  level: "INFO"
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
