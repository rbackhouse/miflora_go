package sensor

import (
	"encoding/binary"
	"time"

	logger "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

type SensorReadings struct {
	Temperature uint16
	Lux         uint32
	Moisture    int
	Fertility   uint16
}

var (
	adapter                      = bluetooth.DefaultAdapter
	SERVICE_UUID                 = bluetooth.NewUUID([16]byte{0x00, 0x00, 0x12, 0x04, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb})
	DATA_CHARACTERISTIC_UUID     = bluetooth.NewUUID([16]byte{0x00, 0x00, 0x1a, 0x01, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb})
	FIRMWARE_CHARACTERISTIC_UUID = bluetooth.NewUUID([16]byte{0x00, 0x00, 0x1a, 0x02, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb})
	REALTIME_CHARACTERISTIC_UUID = bluetooth.NewUUID([16]byte{0x00, 0x00, 0x1a, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb})
)

func connect(deviceAddr bluetooth.Address) ([]bluetooth.DeviceCharacteristic, error) {
	//logger.Infof("Connecting to %s ...", deviceAddr.String())
	device, err := adapter.Connect(deviceAddr, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(30 * time.Second),
	})
	if err != nil {
		logger.Warnf("Failed to connect: %s", err.Error())
		return nil, err
	}
	//logger.Infof("Connected to %s %s ...", deviceAddr.String(), device.Address.String())

	services, err := device.DiscoverServices([]bluetooth.UUID{SERVICE_UUID})
	if err != nil {
		logger.Warnf("Failed to discover the Service: %s", err.Error())
		return nil, err
	}
	service := services[0]

	characteristics, err := service.DiscoverCharacteristics([]bluetooth.UUID{DATA_CHARACTERISTIC_UUID, FIRMWARE_CHARACTERISTIC_UUID, REALTIME_CHARACTERISTIC_UUID})
	if err != nil {
		logger.Warnf("Failed to discover characteristics: %s", err.Error())
		return nil, err
	}
	return characteristics, nil
}

func getBatteryLevel(deviceAddr bluetooth.Address) (int, string, error) {
	characteristics, err := connect(deviceAddr)
	//defer device.Disconnect()
	if err == nil {
		for _, characteristic := range characteristics {
			if characteristic.UUID() == FIRMWARE_CHARACTERISTIC_UUID {
				buf := make([]byte, 7)
				_, err = characteristic.Read(buf)
				if err != nil {
					logger.Warnf("Failed to read firmware: %s", err.Error())
					return -1, "", err
				}
				firmware := string(buf[2:])
				return int(buf[0]), firmware, nil
			}
		}
	}
	return -1, "", err
}

func GetBatteryLevel(deviceAddr bluetooth.Address) (int, string, error) {
	var err error

	for i := 0; i < 5; i++ {
		batteryLevel, firmware, err := getBatteryLevel(deviceAddr)
		if err == nil {
			return batteryLevel, firmware, nil
		}
		logger.Info("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
	return -1, "", err
}

func getReadings(deviceAddr bluetooth.Address) (SensorReadings, error) {
	characteristics, err := connect(deviceAddr)
	if err == nil {
		var data bluetooth.DeviceCharacteristic
		var realtime bluetooth.DeviceCharacteristic

		for _, characteristic := range characteristics {
			if characteristic.UUID() == DATA_CHARACTERISTIC_UUID {
				data = characteristic
			} else if characteristic.UUID() == REALTIME_CHARACTERISTIC_UUID {
				realtime = characteristic
			}
		}
		out := []byte{0xA0, 0x1F}
		//_, err := realtime.WriteWithoutResponse(out)
		_, err := realtime.Write(out)
		if err != nil {
			logger.Warnf("Failed to write data: %s", err.Error())
			return SensorReadings{}, err
		}

		buf := make([]byte, 16)

		_, err = data.Read(buf)
		//logger.Infof("Data: %v", buf)
		if err != nil {
			logger.Warnf("Failed to read data: %s", err.Error())
			return SensorReadings{}, err
		}
		return SensorReadings{
			Temperature: binary.LittleEndian.Uint16(buf[0:2]) / 10,
			Lux:         binary.LittleEndian.Uint32(buf[3:7]),
			Moisture:    int(buf[7]),
			Fertility:   binary.LittleEndian.Uint16(buf[8:10]),
		}, nil
	} else {
		return SensorReadings{}, err
	}
}

func GetReadings(deviceAddr bluetooth.Address) (SensorReadings, error) {
	var err error

	for i := 0; i < 5; i++ {
		sensorReadings, err := getReadings(deviceAddr)
		if err == nil {
			return sensorReadings, nil
		}
		logger.Info("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
	return SensorReadings{}, err
}
