package publisher

import (
	"context"
	"strconv"

	"github.com/masanetes/loudspeaker-runtime/pkg/utils"
	loudspeakerv1alpha1 "github.com/masanetes/loudspeaker/api/v1alpha1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/watch"
)

type Publisher interface {
	Add(obj interface{})
	Update(obj interface{})
	Run(ctx context.Context)
}

type publisher struct {
	initCh    chan *[]loudspeakerv1alpha1.Subscribe
	updateCh  chan *[]loudspeakerv1alpha1.Subscribe
	clientset *kubernetes.Clientset
	canceler  context.CancelFunc
}

func New(clientset *kubernetes.Clientset) Publisher {
	addCh := make(chan *[]loudspeakerv1alpha1.Subscribe)
	updateCh := make(chan *[]loudspeakerv1alpha1.Subscribe)
	p := publisher{
		initCh:    addCh,
		updateCh:  updateCh,
		clientset: clientset,
	}
	return &p
}

func (p *publisher) Add(obj interface{}) {
	config := utils.ConfigDecode(obj)
	p.initCh <- config
	log.Infof("Initialize the capture settings: %+v", *config)
}

func (p *publisher) Update(obj interface{}) {
	config := utils.ConfigDecode(obj)
	p.updateCh <- config
	log.Infof("Configmap has been updated to reflect the capture settings: %+v", *config)
}

func (p *publisher) Run(ctx context.Context) {
	log.Info("Starting subscribe event process...")
	p.processManager(ctx)
}

func (p *publisher) processManager(ctx context.Context) {
	go func() {
		for configs := range p.initCh {
			c, cancel := context.WithCancel(ctx)
			p.canceler = cancel
			p.runner(c, configs)
		}
	}()

	go func() {
		for configs := range p.updateCh {
			p.canceler()
			c, cancel := context.WithCancel(ctx)
			p.canceler = cancel
			p.runner(c, configs)
		}
	}()
	<-ctx.Done()
	return
}

func (p *publisher) runner(ctx context.Context, subscribes *[]loudspeakerv1alpha1.Subscribe) {
	for _, s := range *subscribes {
		go p.process(ctx, s)
	}
}

func (p *publisher) process(ctx context.Context, subscribe loudspeakerv1alpha1.Subscribe) {
	rw, err := p.watcher(subscribe, ctx)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("[Capture Namespace: %s] start", subscribe.Namespace)
	for {
		select {
		case event := <-rw.ResultChan():
			if p, ok := event.Object.(*v1.Event); ok {
				if ok {
					log.Infof("[Capture thread: %s] Type: %s, Namespace: %s, Message: %s",
						subscribe.Namespace,
						p.Type,
						p.Namespace,
						p.Message)
				}
			}
		case <-ctx.Done():
			log.Infof("[Capture Namespace: %s] ended", subscribe.Namespace)
			return
		}
	}
}

func (p *publisher) watcher(subscribe loudspeakerv1alpha1.Subscribe, ctx context.Context) (*watch.RetryWatcher, error) {
	lastResourceVersion := 1
	eventList, err := p.clientset.CoreV1().Events(subscribe.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	for _, p := range eventList.Items {
		if v, _ := strconv.Atoi(p.ObjectMeta.GetResourceVersion()); v > lastResourceVersion {
			lastResourceVersion = v
		}
	}
	return watch.NewRetryWatcher(
		strconv.Itoa(lastResourceVersion),
		cache.NewListWatchFromClient(
			p.clientset.CoreV1().RESTClient(),
			"events",
			subscribe.Namespace,
			fields.Everything(),
		),
	)
}
