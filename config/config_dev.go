package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/masanetes/loudspeaker/api/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type Config struct {
	Type          v1alpha1.ListenerType `envconfig:"TYPE" default:"sentry"`
	ConfigmapName string                `envconfig:"CONFIGMAP" default:"loudspeaker-sample-bar"`
}

func (c *Config) LoadEnv() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}
	return nil
}

func LoadClusterConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/masa/.kube/config")
	if err != nil {
		log.Fatal(err)
	}
	return config
}
