package scanner

import (
	"os"
	"os/signal"
	"syscall"

	"potpie.org/miflora/src/config"
	d "potpie.org/miflora/src/device"
	"potpie.org/miflora/src/handler"

	logger "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func Scan(cfg config.Config) {
	readyChan := make(chan struct{})

	deviceMap := make(map[string]d.DeviceHandler)

	for _, sensorConf := range cfg.Sensors {
		logger.Infof("Loaded config for %s", sensorConf.DeviceId)
		deviceMap[sensorConf.DeviceId] = handler.NewHandler(cfg, sensorConf)
	}

	must("enable BLE stack", adapter.Enable())

	go func() {
		logger.Info("Starting Scan for MiFlora Devices....")
		adapter.Scan(func(adapter *bluetooth.Adapter, scanResult bluetooth.ScanResult) {
			dev, ok := deviceMap[scanResult.Address.String()]
			if ok {
				dev.SetLoaded()
				dev.SetAddress(scanResult.Address)
			}
			stopscan := true
			for _, dev := range deviceMap {
				if !dev.IsLoaded() {
					stopscan = false
				}
			}
			if stopscan {
				logger.Info("Stopping Scan for BLE Devices....")
				adapter.StopScan()
				readyChan <- struct{}{}
			}
		})
	}()

	for {
		select {
		case <-readyChan:
			for _, dev := range deviceMap {
				dev.Handle()
			}
			break
		}
	}
}

func Discover() {
	doneChan := make(chan struct{})

	must("enable BLE stack", adapter.Enable())

	configMap := make(map[string]config.SensorConfig)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		sensors := make([]config.SensorConfig, 0, len(configMap))
		for _, sensor := range configMap {
			sensors = append(sensors, sensor)
		}
		logger.Infof("Saving Discovered Sensors to config file")
		config.SaveConfig(config.Config{Sensors: sensors})
		os.Exit(1)
	}()

	go func() {
		logger.Info("Starting Discovery of MiFlora Devices....")
		adapter.Scan(func(adapter *bluetooth.Adapter, scanResult bluetooth.ScanResult) {
			if scanResult.LocalName() == "Flower care" {
				_, ok := configMap[scanResult.Address.String()]
				if !ok {
					logger.Infof("Found MiFlora sensor : %s %d %s", scanResult.Address.String(), scanResult.RSSI, scanResult.LocalName())
					configMap[scanResult.Address.String()] = config.SensorConfig{
						DeviceId:             scanResult.Address.String(),
						Name:                 "Sensor-" + scanResult.Address.String(),
						ReadingsInterval:     "1m",
						BatteryLevelInterval: "12h",
						MoistureMax:          0,
						MoistureMin:          0,
					}
				}
			}
		})
	}()

	for {
		select {
		case <-doneChan:
			break
		}
	}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
