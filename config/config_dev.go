package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func LoadClusterConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/masa/.kube/config")
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
