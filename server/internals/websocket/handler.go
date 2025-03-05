package websocket

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handles connection with client and python server
func handleConnections(w http.ResponseWriter, r *http.Request, listeningConnection net.Conn) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket: ", err)
		return
	}
	log.Println("Upgraded HTTP connection to WebSocket")
	defer ws.Close()

	// persistent listening connection
	go sendToClient(ws, listeningConnection)

	// temporary sending connection
	sendToPythonServer(ws)
}

// forwards messages from the python server to the client
func sendToClient(ws *websocket.Conn, listeningConnection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := listeningConnection.Read(buffer)
		if err != nil {
			log.Println("Error reading from Python server: ", err)
			return
		}
		log.Println("Received message from Python server: ", string(buffer[:n]))

		if err := ws.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
			log.Println("Error writing to WebSocket:", err)
			return
		}
		log.Println("Sent message to WebSocket client: ", string(buffer[:n]))
	}
}

// reads message from client and handles request
func sendToPythonServer(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading from WebSocket client: ", err)
			return
		}
		log.Println("Received message from WebSocket client: ", string(message))

		handleRequest(ws, message)
	}
}

// forwards message to python server and sends response to client
func handleRequest(ws *websocket.Conn, message []byte) {
	sendingConnection, err := CreateSendingConnection()
	if err != nil {
		log.Println("Error creating sending connection: ", err)
		return
	}
	defer sendingConnection.Close()

	if _, err := sendingConnection.Write(message); err != nil {
		log.Println("Error writing to Python server: ", err)
		return
	}
	log.Println("Sent message to Python server: ", string(message))

	response := make([]byte, 1024)
	n, err := sendingConnection.Read(response)
	if err != nil {
		log.Println("Error reading from Python server: ", err)
		return
	}
	log.Println("Received response from Python server: ", string(response[:n]))

	if err := ws.WriteMessage(websocket.TextMessage, response[:n]); err != nil {
		log.Println("Error writing to WebSocket: ", err)
		return
	}
	log.Println("Sent response to WebSocket client: ", string(response[:n]))
}

// starts the websocket server
func StartWebsocketServer(listeningConnection net.Conn) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(w, r, listeningConnection)
	})

	log.Println("Starting WebSocket server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting WebSocket server: ", err)
	}
}
