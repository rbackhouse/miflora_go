package device

import (
	"potpie.org/miflora/src/sensor"
	"tinygo.org/x/bluetooth"
)

type Device struct {
	Address           string
	RSSI              int16
	Timestamp         int64
	ManufacturingData map[uint16]interface{}
}

type DeviceHandler interface {
	IsLoaded() bool
	SetLoaded()
	SetAddress(addr bluetooth.Address)
	Handle()
}

type Handler interface {
	HandleReadings(sensorReadings sensor.SensorReadings, sensorName string, deviceAddress string)
	HandleBatteryLevel(batteryLevel int, sensorName string, deviceAddress string)
}
