package email

import (
	"fmt"

	"potpie.org/miflora/src/config"
	"potpie.org/miflora/src/device"
	"potpie.org/miflora/src/sensor"
)

func NewEmailHandler(cfg config.Config) device.Handler {
	return &email{
		toAddresses: cfg.Email.ToAddresses,
		config:      cfg.Email,
	}
}

func (e *email) HandleReadings(sensorReadings sensor.SensorReadings, sensorName string, deviceAddress string) {
	message := fmt.Sprintf("To: %v\r\nSubject: Sensor Readings for %s\r\n\r\nSensor Readings for %s:\n\nTemperature: %d\nLux: %d\nMoisture: %d", e.toAddresses, sensorName, sensorName, sensorReadings.Temperature, sensorReadings.Lux, sensorReadings.Moisture)
	sendEmail(e.config, message, e.toAddresses)
}

func (e *email) HandleBatteryLevel(batteryLevel int, firmware string, sensorName string, deviceAddress string) {
	message := fmt.Sprintf("To: %v\r\nSubject: Sensor Battery Level for %s\r\n\r\nBattery Level: %d%%\r\nFirmware: %s", e.toAddresses, sensorName, batteryLevel, firmware)
	sendEmail(e.config, message, e.toAddresses)
}
