//go:build prod
// +build prod

package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/ureuzy/loudspeaker-runtime/pkg/constants"
	"github.com/ureuzy/loudspeaker/api/v1alpha1"
	"k8s.io/client-go/rest"
)

type Config struct {
	Type          v1alpha1.ListenerType `envconfig:"TYPE" required:"true"`
	ConfigmapName string                `envconfig:"CONFIGMAP" required:"true"`
	Namespace     string                `envconfig:"NAMESPACE" required:"true"`
}

func (c *Config) LoadEnv() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}
	return nil
}

type SentryCredentials struct {
	Dsn string `yaml:"dsn"`
}

func (s *SentryCredentials) Load() error {
	c, err := os.ReadFile(constants.CredentialsPath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(c, s); err != nil {
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
