package listener

import (
	"encoding/json"
	"fmt"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/ureuzy/loudspeaker-runtime/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

func NewSentryClient() Client {
	return &client{}
}

func (c *client) Send(kubeEvent *corev1.Event) {

	jsonStr, err := json.Marshal(kubeEvent)
	if err != nil {
		log.Error(err)
	}

	var x map[string]interface{}
	err = json.Unmarshal(jsonStr, &x)
	if err != nil {
		log.Error(err)
	}

	result := &map[string]string{}
	utils.Flatten("", x, result)
	delete(*result, "metadata.managedFields.0")

	event := &sentry.Event{
		Message: fmt.Sprintf("%s.%s.%s",
			kubeEvent.APIVersion,
			kubeEvent.InvolvedObject.Kind,
			kubeEvent.InvolvedObject.Name,
		),
		Tags: *result,
	}

	sentry.CaptureEvent(event)
}
