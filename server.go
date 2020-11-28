package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ChannelMap map[*chan []byte]bool

var chanMap = make(ChannelMap)

func main() {
	fmt.Println("Global Chap App 1.0")

	http.Handle("/", http.FileServer(http.Dir("client/")))
	http.HandleFunc("/ws", serveWS)

	const port = "8080"
	log.Printf("Listening on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	// Create connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// First message is name
	msgType, name, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("new user: %s\n", name)

	// Show number of active users
	res := []byte(fmt.Sprintf("Welcome %s, there are %d active users right now\n", name, len(chanMap)))
	if err := conn.WriteMessage(msgType, res); err != nil {
		log.Println(err)
		return
	}

	// Create channel for user
	msg := make(chan []byte)
	chanMap[&msg] = true

	go func() {
		defer delete(chanMap, &msg)

		for {
			_, bytes, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			msg := fmt.Sprintf("%s: %s\n", name, bytes)
			log.Print(msg)

			chanMap.Broadcast(msg)
		}
	}()

	// Writer
	go func() {
		defer delete(chanMap, &msg)

		for {
			if err = conn.WriteMessage(msgType, <-msg); err != nil {
				log.Println(err)
				return
			}
		}
	}()
}

func (c *ChannelMap) Broadcast(msg string) {
	for ch, ok := range chanMap {
		if ok {
			*ch <- []byte(msg)
		}
	}
}
