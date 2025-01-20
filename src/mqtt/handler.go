package mqtt

import (
	"encoding/json"
	"time"

	"potpie.org/miflora/src/config"
	"potpie.org/miflora/src/device"
	"potpie.org/miflora/src/sensor"
)

func NewMqttHandler(cfg config.Config) device.Handler {
	NewMqttClient(cfg.Mqtt)
	return &mqttData{cfg.Mqtt}
}

func (m *mqttData) HandleReadings(sensorReadings sensor.SensorReadings, sensorName string, deviceAddress string) {
	values := map[string]interface{}{
		"temperature": sensorReadings.Temperature,
		"lux":         sensorReadings.Lux,
		"moisture":    sensorReadings.Moisture,
		"fertility":   sensorReadings.Fertility,
	}
	data := &DeviceData{
		DeviceName:    sensorName,
		DeviceAddress: deviceAddress,
		Timestamp:     time.Now().Unix(),
		Values:        values,
	}
	messageJSON, err := json.Marshal(data)
	if err != nil {
	}
	publishMessage(m.config.MqttTopic, messageJSON)
}

func (m *mqttData) HandleBatteryLevel(batteryLevel int, sensorName string, deviceAddress string) {
	values := map[string]interface{}{
		"batteryLevel": batteryLevel,
	}
	data := &DeviceData{
		DeviceName:    sensorName,
		DeviceAddress: deviceAddress,
		Timestamp:     time.Now().Unix(),
		Values:        values,
	}
	messageJSON, err := json.Marshal(data)
	if err != nil {
	}
	publishMessage(m.config.MqttTopic, messageJSON)
}
