package main

import (
	"log"
	"os"
	"os/signal"
	"server/internals/websocket"
	"syscall"
)

func main() {
	listeningConnection, err := websocket.CreateListeningConnection()
	if err != nil {
		log.Fatal("Error creating listening connection: ", err)
	}
	defer listeningConnection.Close()

	go websocket.StartWebsocketServer(listeningConnection)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	<-signalChannel
	log.Println("Shutting down server")
}
