# ezproxy-log-analyzer

Get usage from your EZProxy log files

## Usage

Start an [ezpaarse](https://github.com/ezpaarse-project/ezpaarse) instance

```
git clone https://github.com/ezpaarse-project/ezpaarse`
cd ezpaarse
docker compose pull
docker compose up -d
cd ..
```

Aggregate and process ezproxy log files

```
scp ezproxy.host:/path/to/ezproxy/logs/ .
./process.sh
```

See usage by platform

```
jq -r '.[]| "\(.platform_name) \(.publisher_name)"' logs/ezpaarse.json|sort|uniq -c|sort -n
```
