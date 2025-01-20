package mqtt

import (
	"fmt"

	logger "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"potpie.org/miflora/src/config"
)

type DeviceData struct {
	DeviceName    string                 `json:"deviceName"`
	DeviceAddress string                 `json:"deviceAddress"`
	Timestamp     int64                  `json:"timestamp"`
	Values        map[string]interface{} `json:"values"`
}

type mqttData struct {
	config config.MqttConfig
}

var client mqtt.Client

func NewMqttClient(config config.MqttConfig) mqtt.Client {
	if client == nil {
		opts := mqtt.NewClientOptions()
		brokerStr := fmt.Sprintf("mqtt://%s:%d", config.MqttBroker, config.MqttPort)
		logger.Infof("Connecting to %s mqtt broker", brokerStr)
		opts.AddBroker(brokerStr)
		opts.SetClientID("miflora")
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		logger.Infof("Connected to %s mqtt broker", brokerStr)
	}
	return client
}

func publishMessage(topic string, msg []byte) {
	logger.Infof("Publishing %s to topic %s", string(msg), topic)
	token := client.Publish(topic, 0, false, msg)
	token.Wait()
}
