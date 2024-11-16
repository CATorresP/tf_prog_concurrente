package utilwebsocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWs[T any, R any](w http.ResponseWriter, r *http.Request, handleFunc func(*T, *R) error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: handleWs: Error upgrading to WebSocket: %v\n", err)
		return
	}
	defer ws.Close()

	for {
		var request T
		err := ws.ReadJSON(&request)
		if err != nil {
			log.Printf("ERROR: handleWs: Error reading JSON: %v\n", err)
			break
		}
		var response R
		err = handleFunc(&request, &response)
		if err != nil {
			log.Printf("ERROR: handleWs: Error handling request: %v\n", err)
			break
		}
		err = ws.WriteJSON(response)
		if err != nil {
			log.Printf("ERROR: handleWs: Error sending response: %v\n", err)
			break
		}
	}
}

func HandleWsFunc[T any, R any](path string, handleFunc func(*T, *R) error) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { handleWs[T, R](w, r, handleFunc) })
}

func handleWsNoRequest[T any](w http.ResponseWriter, r *http.Request, handleFunc func(*T)) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: handleWs: Error upgrading to WebSocket: %v\n", err)
		return
	}
	defer ws.Close()

	for {
		err = ws.ReadJSON(nil)
		if err != nil {
			log.Printf("ERROR: handleWs: Error reading JSON: %v\n", err)
			break
		}
		var response T
		handleFunc(&response)
		err = ws.WriteJSON(response)
		if err != nil {
			log.Printf("ERROR: handleWs: Error sending response: %v\n", err)
			break
		}
	}
}

func HandleWsFuncNoRequest[T any](path string, handleFunc func(*T)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { handleWsNoRequest[T](w, r, handleFunc) })
}
