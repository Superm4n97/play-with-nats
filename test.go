package main

import (
	"github.com/nats-io/nats.go"
	"k8s.io/klog/v2"
)

const (
	serverAddress = "nats://127.0.0.1:30222"
)

func main() {
	nc, err := nats.Connect(serverAddress, nats.UserCredentials("/home/rasel/Downloads/nats.creds"))
	if err != nil {
		klog.Errorf(err.Error())
		return
	}
	js, err := nc.JetStream()
	if err != nil {
		klog.Errorf(err.Error())
		return
	}

	strInfo, err := js.AddStream(&nats.StreamConfig{
		Name:     "LOG",
		Subjects: []string{"natjobs.resp.>"},
	})
	if err != nil {
		klog.Errorf(err.Error())
		return
	}
	connInfo, err := js.AddConsumer(strInfo.Config.Name, &nats.ConsumerConfig{
		Durable:       "MONITOR",
		AckPolicy:     nats.AckExplicitPolicy,
		FilterSubject: "natjobs.resp.>",
	})
	if err != nil {
		klog.Errorf(err.Error())
		return
	}

	sub, err := js.PullSubscribe(connInfo.Config.FilterSubject, "MONITOR")
	if err != nil {
		klog.Errorf(err.Error())
		return
	}

	for {
		msgs, err := sub.Fetch(1)
		if err != nil {
			klog.Errorf(err.Error())
			break
		}

		klog.Infof("received: %s\n", string(msgs[0].Data))
		err = msgs[0].Ack()
		if err != nil {
			klog.Errorf(err.Error())
			break
		}
	}
}
