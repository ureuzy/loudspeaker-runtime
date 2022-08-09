//go:build !prod
// +build !prod

package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/ureuzy/loudspeaker-runtime/pkg/constants"
	"github.com/ureuzy/loudspeaker/api/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	Type          v1alpha1.ListenerType `envconfig:"TYPE" default:"sentry"`
	ConfigmapName string                `envconfig:"CONFIGMAP" default:"loudspeaker-sample-foo"`
	Namespace     string                `envconfig:"NAMESPACE" default:"default"`
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
	kubeconfig := fmt.Sprintf("%s/.kube/config", homedir.HomeDir())
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
