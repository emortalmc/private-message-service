package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	RabbitMQ            RabbitMQConfig            `yaml:"rabbitmq"`
	RelationshipService RelationshipServiceConfig `yaml:"relationship_service"`
	Development         bool                      `yaml:"debug"`

	Port uint16 `yaml:"port"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RelationshipServiceConfig struct {
	Host string `yaml:"host"`
}

func LoadGlobalConfig() (cfg *Config, err error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("cfg")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return
	}

	return
}
