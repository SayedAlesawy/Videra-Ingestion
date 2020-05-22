package tcp

import "github.com/pebbe/zmq4"

// Connection Represents a TCP Connection object
type Connection struct {
	socket *zmq4.Socket //A TCP socket
}
