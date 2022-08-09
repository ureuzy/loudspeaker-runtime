package main

import (
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/ureuzy/loudspeaker-runtime/config"
	"github.com/ureuzy/loudspeaker-runtime/pkg"
	"github.com/ureuzy/loudspeaker-runtime/pkg/listener"
	"github.com/ureuzy/loudspeaker-runtime/pkg/manager"
	"github.com/ureuzy/loudspeaker-runtime/pkg/signals"
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

	var creds config.SentryCredentials
	if err = creds.Load(); err != nil {
		log.Fatal(err)
	}

	if err = sentry.Init(sentry.ClientOptions{
		Dsn: creds.Dsn,
	}); err != nil {
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
