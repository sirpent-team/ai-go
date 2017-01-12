package sirpent

import (
	"encoding/json"
	"net"
	"time"
)

type jsonSocket struct {
	// Opened socket.
	socket  net.Conn
	timeout time.Duration
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
		timeout: timeout,
		encoder: json.NewEncoder(socket),
		decoder: json.NewDecoder(socket),
	}, err
}

func (js *jsonSocket) sendOrTimeout(m interface{}) error {
	js.socket.SetDeadline(time.Now().Add(js.timeout))
	err := js.encoder.Encode(m)
	js.socket.SetDeadline(time.Time{})
	return err
}

func (js *jsonSocket) receiveOrTimeout(r interface{}) error {
	js.socket.SetDeadline(time.Now().Add(js.timeout))
	err := js.decoder.Decode(r)
	js.socket.SetDeadline(time.Time{})
	return err
}
