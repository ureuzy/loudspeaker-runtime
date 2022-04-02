package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/masanetes/loudspeaker-runtime/config"
	"github.com/masanetes/loudspeaker-runtime/pkg"
	"github.com/masanetes/loudspeaker-runtime/pkg/listener"
	"github.com/masanetes/loudspeaker-runtime/pkg/manager"
	"github.com/masanetes/loudspeaker-runtime/pkg/signals"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
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

	err = sentry.Init(sentry.ClientOptions{
		Dsn: "",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	client := listener.NewSentryClient()

	mgr := manager.New(clientset, client)
	configmapsController := pkg.NewConfigmapsController(clientset, &cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mgr.FetchConfig(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			mgr.FetchConfig(newObj)
		},
	}, conf)

	ctx := signals.NewContext()

	go configmapsController.Run(ctx.Done())
	go mgr.Run(ctx)
	<-ctx.Done()
}
