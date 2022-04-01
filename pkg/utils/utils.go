package utils

import (
	"bytes"
	"errors"
	loudspeakerv1alpha1 "github.com/masanetes/loudspeaker/api/v1alpha1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func typecast(obj interface{}) (*v1.ConfigMap, error) {
	v, ok := obj.(*v1.ConfigMap)
	if !ok {
		return nil, errors.New("failed type cast")
	}
	return v, nil
}

func ConfigDecode(obj interface{}) *[]loudspeakerv1alpha1.Subscribe {
	cm, err := typecast(obj)
	if err != nil {
		log.Error(err)
		return nil
	}
	var subscribeConfig []loudspeakerv1alpha1.Subscribe
	d := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(cm.Data["subscribes"])), 4096)
	err = d.Decode(&subscribeConfig)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &subscribeConfig
}
