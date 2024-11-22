# ezproxy-log-analyzer

Get usage from your EZProxy log files

## Usage

Start an [ezpaarse](https://github.com/ezpaarse-project/ezpaarse) instance

```
git clone https://github.com/ezpaarse-project/ezpaarse
cd ezpaarse
docker compose pull
docker compose up -d
cd ..
```

Aggregate and process ezproxy log files

```
scp -r ezproxy.host:/path/to/ezproxy/logs .
./process.sh
```

Process a clean JSON file you can upload into an analytics tool

```
go run main.go
```
