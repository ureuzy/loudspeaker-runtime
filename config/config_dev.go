//go:build !prod
// +build !prod

package config

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/masanetes/loudspeaker/api/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
