package config

import (
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	HandlerType string
	Email       EmailConfig
	Http        HttpConfig
	Mqtt        MqttConfig
	Sensors     []SensorConfig
}

type EmailConfig struct {
	SmtpHost    string
	SmtpPort    int
	Password    string
	FromAddress string
	ToAddresses []string
}

type HttpConfig struct {
	HttpHost  string
	HttpPort  int
	URLSuffix string
}

type MqttConfig struct {
	MqttBroker string
	MqttPort   int
	MqttTopic  string
}

type SensorConfig struct {
	DeviceId             string
	Name                 string
	ReadingsInterval     string
	BatteryLevelInterval string
	MoistureMax          int
	MoistureMin          int
}

func NewConfig() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var config Config

	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)

	if err != nil {
		logger.Warnf("Unable to decode into struct, %v", err)
	}

	return config
}

func SaveConfig(config Config) {
	viper.SetConfigType("yml")
	viper.Set("sensors", config.Sensors)
	viper.WriteConfigAs("discovered.yaml")
}
