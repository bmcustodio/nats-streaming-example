package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

var (
	clientID    string
	clusterID   string
	durableName string
	natsURL     string
	queueGroup  string
	subject     string
)

func main() {
	// Parse command-line flags.
	flag.StringVar(&clientID, "client-id", "", "the nats streaming client id to use")
	flag.StringVar(&clusterID, "cluster-id", "", "the id of the nats streaming cluster to connect to")
	flag.StringVar(&durableName, "durable-name", "", "the name of the durable to use")
	flag.StringVar(&natsURL, "nats-url", "", "the url at which the nats server can be reached")
	flag.StringVar(&queueGroup, "queue-group", "", "the queue group to use")
	flag.StringVar(&subject, "subject", "date", "the subject to subscribe to")
	flag.Parse()

	// Configure logging.
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	// Connect to the NATS Streaming server.
	c, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("failed to connect to the nats streaming server: %v", err)
	}
	log.Infof("connected: %t", c.NatsConn().IsConnected())

	// Subscribe to the 'subject' subject using the specified queue group and durable name.
	s, err := c.QueueSubscribe("date", queueGroup, func(m *stan.Msg) {
		log.Infof("message received: %v", m.String())
	}, stan.DurableName(durableName))
	if err != nil {
		log.Fatalf("failed to create subscription: %v", err)

	}

	// Close the subscription and connection when we terminate.
	defer func() {
		if err := s.Close(); err != nil {
			log.Fatalf("failed to close subscription: %v", err)
		}
		if err := c.Close(); err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
		log.Info("bye")
	}()

	// Wait for SIGINT.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
}
