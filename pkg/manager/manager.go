package manager

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/ureuzy/loudspeaker-runtime/pkg/listener"
	"github.com/ureuzy/loudspeaker-runtime/pkg/utils"
	loudspeakerv1alpha1 "github.com/ureuzy/loudspeaker/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/watch"
)

type Manager interface {
	FetchConfig(obj interface{})
	Run(ctx context.Context)
}

type manager struct {
	configCh  chan *[]loudspeakerv1alpha1.Observe
	clientset *kubernetes.Clientset
	canceler  context.CancelFunc
	listener  listener.Client
}

func New(clientset *kubernetes.Clientset, client listener.Client) Manager {
	configCh := make(chan *[]loudspeakerv1alpha1.Observe)
	p := manager{
		configCh:  configCh,
		clientset: clientset,
		canceler:  nil,
		listener:  client,
	}
	return &p
}

func (p *manager) FetchConfig(obj interface{}) {
	config := utils.ConfigDecode(obj)
	p.configCh <- config
	log.Infof("fetch settings: %+v", *config)
}

func (p *manager) Run(ctx context.Context) {
	log.Info("Starting observe event process...")
	for configs := range p.configCh {
		if p.canceler != nil {
			p.canceler()
		}
		p.runner(ctx, configs)
	}
	<-ctx.Done()
	return
}

func (p *manager) runner(ctx context.Context, observes *[]loudspeakerv1alpha1.Observe) {
	c, cancel := context.WithCancel(ctx)
	p.canceler = cancel
	for _, s := range *observes {
		go p.process(c, s)
	}
}

func (p *manager) process(ctx context.Context, observe loudspeakerv1alpha1.Observe) {

	rw, err := p.watcher(observe, ctx)
	if err != nil {
		log.Error(err)
		return
	}

	namespace := observe.Namespace
	if namespace == "" {
		namespace = "*"
	}

	log.Warningf("[%s] Start observation", namespace)
	for {
		select {
		case event := <-rw.ResultChan():
			if e, ok := event.Object.(*v1.Event); ok {
				if observe.IgnoreReasons.Contains(e.Reason) ||
					observe.IgnoreObjectKinds.Contains(e.InvolvedObject.Kind) ||
					observe.IgnoreObjectNames.Contains(e.InvolvedObject.Name) ||
					observe.IgnoreEventTypes.Contains(e.Type) {
					continue
				}
				p.listener.Send(e)
				log.Infof("[%s] %s %s %s", namespace, e.Namespace, e.Name, e.Reason)
			}
		case <-ctx.Done():
			log.Warningf("[%s] Observation is terminated", namespace)
			return
		}
	}
}

func (p *manager) watcher(observe loudspeakerv1alpha1.Observe, ctx context.Context) (*watch.RetryWatcher, error) {
	lastResourceVersion := 1
	eventList, err := p.clientset.CoreV1().Events(observe.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	for _, item := range eventList.Items {
		if v, _ := strconv.Atoi(item.ObjectMeta.GetResourceVersion()); v > lastResourceVersion {
			lastResourceVersion = v
		}
	}
	return watch.NewRetryWatcher(
		strconv.Itoa(lastResourceVersion),
		cache.NewListWatchFromClient(
			p.clientset.CoreV1().RESTClient(),
			"events",
			observe.Namespace,
			fields.Everything(),
		),
	)
}
