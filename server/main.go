package main

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

func (ws *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Unable to upgrade WebSocket connection due: \n\t[%v]", err)
		return
	}

	defer func() {
		log.Warn("Closing current connection.")
		c.Close()
	}()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Errorf("Unable to reade message due: \n\t[%v]", err)
			return
		}
		if mt == websocket.BinaryMessage {
			err := c.WriteMessage(websocket.TextMessage, []byte("Server does not allow binary."))
			if err != nil {
				log.Errorf("Unable to send message due: \n\t[%v]", err)
				return
			}
			log.Warnf("No response to binary message.")
			return
		}

		log.Infof("Received the current message: \n\t[%v]", strings.Trim(string(msg), "\n"))
	}

}

func main() {
	wsHandler := WebSocketHandler{
		upgrader: websocket.Upgrader{},
	}
	http.Handle("/", &wsHandler)
	log.Infof("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
