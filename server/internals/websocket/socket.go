package websocket

import (
	"errors"
	"log"
	"net"
	"strings"
)

// performs a handshake with the python server
func handshake(connection net.Conn) error {
	message := []byte("[syn]\x00")
	_, err := connection.Write(message)
	if err != nil {
		connection.Close()
		return err
	}
	log.Println("Handshake message sent to Python server: ", string(message))

	response := make([]byte, 1024)
	n, err := connection.Read(response)
	if err != nil {
		connection.Close()
		return err
	}
	log.Println("Handshake response received from Python server: ", string(response))

	responseStr := strings.Trim(string(response[:n]), "\x00")
	if responseStr != "[ack]" {
		connection.Close()
		return errors.New("Handshake failed")
	}
	log.Println("Handshake successful")

	return nil
}

// creates a listening connection to python server
func CreateListeningConnection() (net.Conn, error) {
	connection, err := net.Dial("tcp", "localhost:5001")
	if err != nil {
		return nil, err
	}

	if err = handshake(connection); err != nil {
		connection.Close()
		return nil, err
	}
	log.Println("Listening connection established with Python server")

	return connection, nil
}

// creates a sending connection to python server
func CreateSendingConnection() (net.Conn, error) {
	connection, err := net.Dial("tcp", "localhost:5001")
	if err != nil {
		return nil, err
	}
	log.Println("Sending connection established with Python server")

	return connection, nil
}
