package pkg

import (
	"github.com/masanetes/loudspeaker-runtime/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func NewConfigmapsController(clientset *kubernetes.Clientset, funcs *cache.ResourceEventHandlerFuncs, config config.Config) cache.Controller {
	watchConfigmaps := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		string(v1.ResourceConfigMaps),
		config.Namespace,
		fields.SelectorFromSet(fields.Set{"metadata.name": config.ConfigmapName}))

	_, configmapsController := cache.NewInformer(watchConfigmaps, &v1.ConfigMap{}, 0, *funcs)

	return configmapsController
}
