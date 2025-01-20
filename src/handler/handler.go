package handler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	logger "github.com/sirupsen/logrus"
	"potpie.org/miflora/src/config"
	"potpie.org/miflora/src/device"
	"potpie.org/miflora/src/email"
	"potpie.org/miflora/src/http"
	"potpie.org/miflora/src/mqtt"
	"potpie.org/miflora/src/sensor"
	"tinygo.org/x/bluetooth"
)

var mu = sync.Mutex{}

type handle struct {
	loaded    bool
	scheduled bool
	config    config.SensorConfig
	addr      bluetooth.Address
	handler   device.Handler
	scheduler gocron.Scheduler
}

func (c *handle) getReadings() {
	mu.Lock()
	defer mu.Unlock()
	//logger.Infof("Getting Sensor Readings for %s:%s", c.config.Name, c.config.DeviceId)
	sensorReadings, err := sensor.GetReadings(c.addr)
	if err == nil {
		logger.Infof("Sensor Readings for %s : %+v", c.config.Name, sensorReadings)

		report := false
		if c.config.MoistureMax == -1 || sensorReadings.Moisture > c.config.MoistureMax {
			report = true
		} else if c.config.MoistureMin == -1 || sensorReadings.Moisture < c.config.MoistureMin {
			report = true
		}
		if report {
			c.handler.HandleReadings(sensorReadings, c.config.Name, c.config.DeviceId)
		}
	}
}

func (c *handle) getBatteryLevel() {
	mu.Lock()
	defer mu.Unlock()
	//logger.Infof("Checking Battery Level on %s", c.config.DeviceId)
	batteryLevel, err := sensor.GetBatteryLevel(c.addr)
	if err == nil {
		logger.Infof("Battery Level for %s : %d\n", c.config.Name, batteryLevel)
		c.handler.HandleBatteryLevel(batteryLevel, c.config.Name, c.config.DeviceId)
	}
}

func (c *handle) Handle() {
	if !c.scheduled {
		c.scheduled = true
		c.scheduler, _ = gocron.NewScheduler()
		readingsInterval, err := time.ParseDuration(c.config.ReadingsInterval)
		if err != nil {
			logger.Error("Error parsing readings interval")
			return
		}
		batteryLevelInterval, err := time.ParseDuration(c.config.BatteryLevelInterval)
		if err != nil {
			logger.Error("Error parsing battery level interval")
			return
		}
		logger.Infof("%s Readings Interval is %s\n", c.config.Name, readingsInterval)
		logger.Infof("%s Battery Level Interval is %s\n", c.config.Name, batteryLevelInterval)

		c.scheduler.NewJob(
			gocron.DurationJob(
				readingsInterval,
			),
			gocron.NewTask(
				func() {
					c.getReadings()
				},
			),
		)
		c.scheduler.NewJob(
			gocron.DurationJob(
				batteryLevelInterval,
			),
			gocron.NewTask(
				func() {
					c.getBatteryLevel()
				},
			),
		)
		c.scheduler.Start()
	}
}

func (c *handle) IsLoaded() bool {
	return c.loaded
}

func (c *handle) SetLoaded() {
	c.loaded = true
}

func (c *handle) SetAddress(addr bluetooth.Address) {
	c.addr = addr
}

func NewHandler(cfg config.Config, sensorCfg config.SensorConfig) device.DeviceHandler {
	var handler device.Handler

	if cfg.HandlerType == "email" {
		handler = email.NewEmailHandler(cfg)
	} else if cfg.HandlerType == "http" {
		handler = http.NewHttpHandler(cfg)
	} else if cfg.HandlerType == "mqtt" {
		handler = mqtt.NewMqttHandler(cfg)
	}

	return &handle{config: sensorCfg, scheduled: false, loaded: false, handler: handler}
}
