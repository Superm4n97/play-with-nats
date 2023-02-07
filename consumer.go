package main

import (
	"github.com/nats-io/nats.go"
	_ "github.com/nats-io/nats.go"
	"k8s.io/klog/v2"
	"log"
	"os"
)

const (
	//natsAddress = "localhost:30222"
	//natsServerAddress = "nats://127.0.0.1:30222"
	// natsServerAddress = "this-is-nats.appscode.ninja:4222"
	//natsServerAddress = "nats://nats.appscode.ninja:4222"

	natsSubject = "stackscript-log"
)

func main() {
	natsServerAddress, ok := os.LookupEnv("SERVER")
	if !ok {
		log.Fatal("missing nats server address")
		return
	}

	var nc *nats.Conn
	var err error

	_, ok = os.LookupEnv("CRED")
	if ok {
		//nc, err := nats.Connect(natsServerAddress, nats.UserCredentials("/home/rasel/Desktop/nats/admin.creds"))
		nc, err = nats.Connect(natsServerAddress, nats.UserCredentials(os.Getenv("CRED")))
		if err != nil {
			klog.Infof("failed to connect with nats server, %s", err.Error())
			return
		}
	} else {
		nc, err = nats.Connect(natsServerAddress, nats.UserInfo(os.Getenv("USER"), os.Getenv("PASS")))
		if err != nil {
			klog.Infof("failed to connect with nats server, %s", err.Error())
		}
	}

	// returns a jetstream context which will be used for message passing
	js, err := nc.JetStream()
	if err != nil {
		klog.Errorf("failed to create Jetstream contex")
		return
	}

	strm, err := addStream(js)
	if err != nil {
		klog.Errorf(err.Error())
		return
	}

	err = addConsumer(js, strm)
	if err != nil {
		klog.Errorf(err.Error())
		return
	}
}

func addStream(js nats.JetStreamContext) (*nats.StreamInfo, error) {

	strInfo, err := js.AddStream(&nats.StreamConfig{
		Name:     "LOG",
		Subjects: []string{natsSubject},
	})
	if err != nil {
		return nil, err
	}
	return strInfo, nil
}

func addConsumer(js nats.JetStreamContext, streamInfo *nats.StreamInfo) error {

	//regular consumer remember their position while they are connected with the client.
	//if the connection is lost, their position will also be lost, Durable remembers their position if the connection is lost
	//Durable subscription identify themselves with a name, connect and disconnect will not affect the durable subscriptions position in the channel
	connInfo, err := js.AddConsumer(streamInfo.Config.Name, &nats.ConsumerConfig{
		Durable:       "MONITOR",
		AckPolicy:     nats.AckExplicitPolicy,
		FilterSubject: natsSubject,
	})
	if err != nil {
		return err
	}

	sub, err := js.PullSubscribe(connInfo.Config.FilterSubject, connInfo.Name, nats.BindStream(connInfo.Stream))
	if err != nil {
		return err
	}
	defer func() {
		err = sub.Unsubscribe()
		if err != nil {
			klog.Infof("failed to unsubscribe the consumer")
		}
	}()

	klog.Infof("waiting for message:")
	for {
		msgs, err := sub.Fetch(1)
		if err == nats.ErrTimeout {
			continue
		}
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			continue
		}

		err = msgs[0].Ack()
		if err != nil {
			return err
		}

		klog.Info("R:", string(msgs[0].Data))
	}
}
