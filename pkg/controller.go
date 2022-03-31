package pkg

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func NewConfigmapsController(clientset *kubernetes.Clientset, funcs * cache.ResourceEventHandlerFuncs) cache.Controller {
	watchConfigmaps := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		string(v1.ResourceConfigMaps),
		"default",
		fields.SelectorFromSet(fields.Set{"metadata.name": "loudspeaker-sample-bar"}))

	_, configmapsController := cache.NewInformer(watchConfigmaps, &v1.ConfigMap{}, 0, *funcs)

	return configmapsController
}
