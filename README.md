# sink_influxdb

`sink_influxdb` provides a simple sink for [statsite](https://github.com/armon/statsite) using [influxdb](https://github.com/influxdb/influxdb) as the backend. 

Currently this is a WIP.

Start InfluxDB
```
docker run -p "8083:8083" -p "8086:8086" -it influxdb
```

Export config ENVs
```
export INFLUX_HOST="127.0.0.1"
export INFLUX_PORT=8086
export INFLUX_DATABASE="test"
export INFLUX_USER="root"
export INFLUX_PASSWORD="test"
export INFLUX_BATCHSIZE=5
export INFLUX_RP="default"
```

Test reading from stdin
```
cat test_data |./bin/sink_influxdb
```

## Build

The project is built using [gb](http://getgb.io) 
```
gb build
```
## Test
```
gb test
```
