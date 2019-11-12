package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CESARBR/knot-babeltower/internal/config"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/server"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

func monitorSignals(sigs chan os.Signal, logger logging.Logger) {
	signal := <-sigs
	logger.Infof("Signal %s received", signal)
}

func main() {
	config := config.Load()
	logrus := logging.NewLogrus(config.Logger.Level)

	logger := logrus.Get("Main")
	logger.Info("Starting KNoT Babeltower")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go monitorSignals(sigs, logger)

	amqp := network.NewAmqpHandler(config.RabbitMQ.URL, logrus.Get("AmqpHandler"))
	amqp.Start()

	server := server.NewServer(config.Server.Port, logrus.Get("Server"))
	server.Start()
}
