package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

var (
	clientID  string
	clusterID string
	natsURL   string
	subject   string
)

func main() {
	// Configure logging.
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	// Parse command-line flags.
	flag.StringVar(&clientID, "client-id", "", "the nats streaming client id to use")
	flag.StringVar(&clusterID, "cluster-id", "", "the id of the nats streaming cluster to connect to")
	flag.StringVar(&natsURL, "nats-url", "", "the url at which the nats server can be reached")
	flag.StringVar(&subject, "subject", "date", "the subject on which to publish")
	flag.Parse()

	// Connect to the NATS Streaming server.
	c, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("failed to connect to the nats streaming server: %v", err)
	}
	log.Infof("connected: %t", c.NatsConn().IsConnected())

	// Publish the current timestamp to the 'subject' subject every second.
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	go func() {
		for {
			<-t.C
			m := time.Now().String()
			log.Infof("publishing message: %v", m)
			if err := c.Publish(subject, []byte(m)); err != nil {
				log.Errorf("failed to publish message: %v", err)
			}
		}
	}()

	// Close the subscription and connection when we terminate.
	defer func() {
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
