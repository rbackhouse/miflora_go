# Xiaomi MiFlora Sensor BLE Application for OSX
Golang based app that can discover and scan sensor values from Xiaomi MiFlora Sensors. While it is using a common golang BLE library from tinygo it is written work correcty on OSX based machines.

## Prereqs
Run go mod tidy to load dependencies

## Discover available sensors
```
go run src/main/main.go discover
```
After ctrl-c is pressed results are found in discovered.yaml

## Scan for readings and battery level
Create a config.yaml file using the discovered MiFlora device UUIDs

```
sensors:
  - deviceId: ""
    name: "Sensor 1"
    readingsInterval: "5m"
    batteryLevelInterval: "1h"
    moistureMax: 0
    moistureMin: 0
  - deviceId: ""
    name: "Sensor 2"
    readingsInterval: "5m"
    batteryLevelInterval: "1h"
    moistureMax: 0
    moistureMin: 0
```
Also see template-config.yaml and set preferred reporting option (email, http or mqtt)

### Setting Moisture Level reporting
If you want to report regardless of the moisture value set moistureMax and moistureMin to -1. Setting them to 0 will ensure that the reporting only occurs when the recorded moisture value is greater that 0 (or the value you have configured)

### Run scan to start recording from configured sensors
```
go run src/main/main.go scan
```

## Build binary

```
go build -o ./miflora -ldflags="-w -s" -v ./src/main/main.go
```

### Running via the binary
```
./miflora scan
```