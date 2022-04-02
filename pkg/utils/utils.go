package utils

import (
	"bytes"
	"errors"
	"fmt"
	loudspeakerv1alpha1 "github.com/masanetes/loudspeaker/api/v1alpha1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"strconv"
)

func typecast(obj interface{}) (*v1.ConfigMap, error) {
	v, ok := obj.(*v1.ConfigMap)
	if !ok {
		return nil, errors.New("failed type cast")
	}
	return v, nil
}

func ConfigDecode(obj interface{}) *[]loudspeakerv1alpha1.Observe {
	cm, err := typecast(obj)
	if err != nil {
		log.Error(err)
		return nil
	}
	var observeConfig []loudspeakerv1alpha1.Observe
	d := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(cm.Data["observes"])), 4096)
	err = d.Decode(&observeConfig)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &observeConfig
}

func Flatten(prefix string, src map[string]interface{}, dest *map[string]string) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			Flatten(prefix+k, child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				(*dest)[prefix+k+"."+strconv.Itoa(i)] = fmt.Sprintf("%v", child[i])
			}
		default:
			if v2 := fmt.Sprintf("%v", v); len(v2) > 0 {
				(*dest)[prefix+k] = v2
			}
		}
	}
}
