package main

import (
	"fmt"
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	fmt.Println("Starting server... TCP")

	server := mqtt.New()
	tcp := listeners.NewTCP("t1", ":1883")

	// auth controller
	// add auth controller
	// server.AddController(&Authenticator)
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(Authenticator),
	})
	if err != nil {
		log.Fatal(err)
	}
	// Start the server
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Started!")
	<-done
	fmt.Println("Caught Signal")
	server.Close()
	fmt.Println("  Finished  ")
}
