# Culvert

A tool to quickly open ssh forwarding tunnels.

## Build

```
git clone https://gitee.com/autom-studio/culvert.git -b static
# modify Version, ConfigYaml, KeyStr in internal/tunnel/config.go
# build
go build -o examples/culvert_<tunnel_name>
```

## Start

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
