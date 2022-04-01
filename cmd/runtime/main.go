package main

import (
	"github.com/masanetes/loudspeaker-runtime/config"
	"github.com/masanetes/loudspeaker-runtime/pkg"
	"github.com/masanetes/loudspeaker-runtime/pkg/publisher"
	"github.com/masanetes/loudspeaker-runtime/pkg/signals"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func main() {

	credentials := config.LoadClusterConfig()
	clientset, err := kubernetes.NewForConfig(credentials)
	if err != nil {
		log.Fatal(err)
	}

	var conf config.Config
	if err = conf.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	pub := publisher.New(clientset)
	configmapsController := pkg.NewConfigmapsController(clientset, &cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pub.Add(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pub.Update(newObj)
		},
	}, conf)

	ctx := signals.NewContext()

	go configmapsController.Run(ctx.Done())
	go pub.Run(ctx)
	<-ctx.Done()
}
