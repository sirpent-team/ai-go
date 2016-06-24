package sirpent

import (
	"encoding/json"
	"net"
	"time"
)

type jsonSocket struct {
	// Opened socket.
	socket  net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
}

func newJsonSocket(server_address string, timeout time.Duration) (*jsonSocket, error) {
	socket, err := net.DialTimeout("tcp", server_address, timeout)
	if err != nil {
		return nil, err
	}

	return &jsonSocket{
		socket:  socket,
		encoder: json.NewEncoder(socket),
		decoder: json.NewDecoder(socket),
	}, err
}

func (js *jsonSocket) sendOrTimeout(m interface{}, timeout time.Duration) error {
	js.socket.SetDeadline(time.Now().Add(timeout))
	err := js.encoder.Encode(m)
	js.socket.SetDeadline(time.Time{})
	return err
}

func (js *jsonSocket) receiveOrTimeout(r interface{}, timeout time.Duration) error {
	js.socket.SetDeadline(time.Now().Add(timeout))
	err := js.decoder.Decode(r)
	js.socket.SetDeadline(time.Time{})
	return err
}
