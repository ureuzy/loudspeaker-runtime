// +build prod

package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/masanetes/loudspeaker/api/v1alpha1"
	"k8s.io/client-go/rest"
)

type Config struct {
	Type          v1alpha1.ListenerType `envconfig:"TYPE" required:"true"`
	ConfigmapName string                `envconfig:"CONFIGMAP" required:"true"`
}

func (c *Config) LoadEnv() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}
	return nil
}

func LoadClusterConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}
