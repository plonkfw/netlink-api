package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	routingv1 "github.com/plonkfw/netlink-api/routing/v1"
	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
)

// Start logging
func init() {
	t := time.Now()
	// WHY isn't this format specification just yyyy-MM-dd etc
	utilsv1.Log.Info().Msg("Logs begin at " + t.UTC().Format("2006-01-02 15:04:05") + " UTC")
}

// Setup http listener and routes
func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for signal := range sigChannel {
			cleanup(signal)
		}
	}()

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "localhost:4829"
	}
	routerv1 := routingv1.NewAPIRouter()
	utilsv1.Log.Fatal().Err(http.ListenAndServe(listen, routerv1))
}

func cleanup(signal os.Signal) {
	msg := fmt.Sprintf("Caught %v signal, exiting", signal)
	utilsv1.Log.Info().Msg(msg)
	t := time.Now()
	// WHY isn't this format specification just yyyy-MM-dd etc
	utilsv1.Log.Info().Msg("Logs end at " + t.UTC().Format("2006-01-02 15:04:05") + " UTC")
	os.Exit(0)
}
