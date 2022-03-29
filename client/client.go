package mainn

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/mochi-co/hanami"
)

const publicBroker = "tcp://localhost:1883"

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- fmt.Errorf("caught signal")
	}()

	// The hanami client takes standard paho options.
	options := paho.NewClientOptions()

	// Create the new hanami client with the broker address and the paho options.
	// 	If you have supplied brokers directly to options, the first field (host:publicBroker)
	// 	will be ignored, so feel free to pass it a "" string, eg. hanami.New("", options)
	client := hanami.New(publicBroker, options)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to", publicBroker)

	cb := func(in *hanami.Payload) {
		log.Printf("RECV: %+v\n", in)
	}

	err = client.Subscribe("example", "hanami/example/+", 0, false, cb)
	if err != nil {
		log.Fatal(err)
	}

	// Post a notice when a message comes in on the hanami/done topic.
	err = client.Subscribe("example", "hanami/done", 0, false, func(in *hanami.Payload) {
		log.Println("Done! You may now ctrl-c...")
	})
	if err != nil {
		log.Fatal(err)
	}

	// This is a second subscription for the topic of hanami/done, which will also be called
	// and run the standard callback, printing the message data.
	err = client.Subscribe("watcher", "hanami/done", 0, false, cb)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish("hanami/example/bool", 0, false, true)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish("hanami/example/num", 0, false, 5)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish("hanami/example/string", 0, false, "testing")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish("hanami/example/map", 1, false, map[string]interface{}{
		"v": "this is my value",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Send a JWT signed message.
	client.Secret = []byte("hanami-test")
	_, err = client.PublishSigned("hanami/example/signed", 1, false, "a signed test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println([]byte("this is my value"))

	_, err = client.Publish("hanami/done", 0, false, true)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for signals...
	<-done
	client.Unsubscribe("example", "hanami/example/+") // Unsubscribe a single filter.
	client.UnsubscribeAll("example", false)           // unsubscribe all filters for the "example" subclient.

}