package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CESARBR/knot-babeltower/internal/config"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/server"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"

	thingDeliveryAmqp "github.com/CESARBR/knot-babeltower/pkg/thing/delivery/amqp"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	msgConsumerAmqp "github.com/CESARBR/knot-babeltower/pkg/thing/handler/amqp"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

func monitorSignals(sigs chan os.Signal, quit chan bool, logger logging.Logger) {
	signal := <-sigs
	logger.Infof("Signal %s received", signal)
	quit <- true
}

func main() {
	config := config.Load()
	logrus := logging.NewLogrus(config.Logger.Level)

	logger := logrus.Get("Main")
	logger.Info("Starting KNoT Babeltower")

	sigs := make(chan os.Signal, 1)
	quit := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go monitorSignals(sigs, quit, logger)

	amqpChan := make(chan bool, 1)
	amqp := network.NewAmqpHandler(config.RabbitMQ.URL, logrus.Get("AmqpHandler"))
	go amqp.Start(amqpChan)

	thingProxy := http.NewThingProxy(logrus.Get("ThingProxy"), config.Things.Hostname, config.Things.Port)

	msgPublisher := thingDeliveryAmqp.NewMsgPublisher(logrus.Get("MsgPublisher"), amqp)
	// registerThing := interactors.NewRegisterThing(logrus.Get("RegisterThing"), msgPublisher, thingProxy)
	// updateSchema := interactors.NewUpdateSchema(logrus.Get("UpdateSchema"), msgPublisher, thingProxy)

	thingInteractor := interactors.NewThingInteractor(logrus.Get("RegisterThing"), msgPublisher, thingProxy)

	msgChan := make(chan bool, 1)
	msgConsumer := msgConsumerAmqp.NewMsgConsumer(logrus.Get("MsgConsumer"), amqp, thingInteractor)

	serverChan := make(chan bool, 1)
	server := server.NewServer(config.Server.Port, logrus.Get("Server"))
	go server.Start(serverChan)

	for {
		select {
		case started := <-serverChan:
			if started {
				logger.Info("Server started")
			}
		case started := <-amqpChan:
			if started {
				logger.Info("AMQP started")
				go msgConsumer.Start(msgChan)
			}
		case started := <-msgChan:
			if started {
				logger.Info("Msg consumer started")
			} else {
				quit <- true
			}
		case <-quit:
			msgConsumer.Stop()
			amqp.Stop()
			server.Stop()
		}
	}
}
