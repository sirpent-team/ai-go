package sirpent

import (
	"golang.org/x/net/websocket"
)

// @TODO: Use field tags to use lowercase JSON names.
// Detailed at https://golang.org/pkg/encoding/json/#Marshal

type API struct {
	Websockets []*websocket.Conn
}

// import (
//   "websocket"
// )

// func (p Player) SendWorld(w *World) error {
//   fmt.Fprintf(Player.Connection, "SENDING_WORLD\r\n")

//   return err

//   //fmt.Fprintf(Player.Connection, "WELCOME TO SIRPENT-GO\r\n\r\n")
//   //status, err := bufio.NewReader(Player.Connection).ReadString('\n')
// }

// // receive JSON type T
// var data T
// websocket.JSON.Receive(ws, &data)

// // send JSON type T
// websocket.JSON.Send(ws, data)
