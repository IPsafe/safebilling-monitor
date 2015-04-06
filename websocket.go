package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/googollee/go-socket.io"
)

func Websocket(wg *sync.WaitGroup, channel chan message) {
	sio := socketio.NewSocketIOServer(&socketio.Config{})

	var msg message
	// Set the on connect handler
	sio.On("connect", func(ns *socketio.NameSpace) {
		log.Println("Connected: ", ns.Id())

		for _, msg := range activeCalls.Calls {
			b, _ := json.Marshal(msg)

			sio.Broadcast("add", string(b))

		}

		sio.Broadcast("connected", ns.Id())

	})

	// Set the on disconnect handler
	sio.On("disconnect", func(ns *socketio.NameSpace) {
		log.Println("Disconnected: ", ns.Id())
		sio.Broadcast("disconnected", ns.Id())
	})

	go func() {
		// Set a handler for news messages
		for {
			msg = <-channel

			b, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch msg.Action {
			case "Add":
				log.Println("Add: ", string(b))
				sio.Broadcast("add", string(b))

			case "Del":
				//	log.Println("Del: ", string(b))
				sio.Broadcast("del", string(b))

			case "Mod":
				//log.Println("Mod: ", string(b))
				sio.Broadcast("mod", string(b))
			}

		}
	}()

	//	sio.On("add", func(ns *socketio.NameSpace, message string) {
	//		log.Println("Add: ", message)
	//	sio.Broadcast("add", message)
	//	})

	// Serve our website
	sio.Handle("/", http.FileServer(http.Dir("./www/")))

	// Start listening for socket.io connections
	println("listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", sio))
	//wg.Done()
	//close(a)

}
