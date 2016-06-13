package main

import (
  "fmt"
  "net/http"
  "golang.org/x/net/websocket"
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
