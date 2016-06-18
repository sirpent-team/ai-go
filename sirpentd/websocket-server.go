package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func worldLiveHandler(ws *websocket.Conn) {
	fmt.Fprintf(ws, "Test %s\r\n", "abc")
}

func main() {
	http.Handle("/worlds/live", websocket.Handler(worldLiveHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
