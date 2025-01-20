package http

import (
	"bytes"
	"encoding/json"

	"potpie.org/miflora/src/config"
	"potpie.org/miflora/src/device"
	"potpie.org/miflora/src/sensor"
)

type SensorData struct {
	SensorName     string                `json:"sensorName"`
	DeviceAddress  string                `json:"deviceAddress"`
	SensorReadings sensor.SensorReadings `json:"sensorReadings"`
}

type BatteryData struct {
	SensorName    string `json:"sensorName"`
	DeviceAddress string `json:"deviceAddress"`
	BatteryLevel  int    `json:"batteryLevel"`
}

func NewHttpHandler(cfg config.Config) device.Handler {
	return &httpData{
		config: cfg.Http,
	}
}

func (h *httpData) HandleReadings(sensorReadings sensor.SensorReadings, sensorName string, deviceAddress string) {
	body := &SensorData{
		SensorName:     sensorName,
		DeviceAddress:  deviceAddress,
		SensorReadings: sensorReadings,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	sendHttpRequest(h.config, payloadBuf, h.config.URLSuffix+"sensordata")
}

func (h *httpData) HandleBatteryLevel(batteryLevel int, sensorName string, deviceAddress string) {
	body := &BatteryData{
		SensorName:    sensorName,
		DeviceAddress: deviceAddress,
		BatteryLevel:  batteryLevel,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	sendHttpRequest(h.config, payloadBuf, h.config.URLSuffix+"sensorbattery")
}
