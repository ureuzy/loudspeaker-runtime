package listener

import corev1 "k8s.io/api/core/v1"

type Client interface {
	Send(event *corev1.Event)
}

type client struct {}
