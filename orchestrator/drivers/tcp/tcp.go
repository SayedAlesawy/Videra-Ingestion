package tcp

import (
	"fmt"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"github.com/pebbe/zmq4"
)

// Connection Represents a TCP Connection object
type Connection struct {
	socket *zmq4.Socket //A TCP socket
}

// NewConnection A function to obtain and initialize a new tcp connection object
func NewConnection(socketType zmq4.Type, topic string) (Connection, bool) {
	socket, err := zmq4.NewSocket(socketType)
	status := errors.IsError(err)

	socket.SetLinger(0)

	if socketType == zmq4.SUB {
		socket.SetSubscribe(topic)
	}

	return Connection{socket: socket}, status
}

// Connect A function that connects a socket to a list of outgoing endpoints
func (connectionObj *Connection) Connect(endpoints ...string) {
	for _, endpoints := range endpoints {
		connectionObj.socket.Connect(endpoints)
	}
}

// Bind A function that binds a socket to a list of incoming endpoints
func (connectionObj *Connection) Bind(endpoints ...string) {
	for _, endpoint := range endpoints {
		connectionObj.socket.Bind(endpoint)
	}
}

// Disconnect A function that terminates connection to a list of endpoints
func (connectionObj *Connection) Disconnect(endpoints ...string) {
	for _, endpoint := range endpoints {
		connectionObj.socket.Disconnect(endpoint)
	}
}

// Close A function to close the connection socket
func (connectionObj *Connection) Close() {
	connectionObj.socket.Close()
}

// Send A function that synchronously sends a msg on the connection
func (connectionObj *Connection) Send(msg interface{}, flags zmq4.Flag) bool {
	var err error

	switch msg.(type) {
	case string:
		_, err = connectionObj.socket.Send(msg.(string), flags)
	case []byte:
		_, err = connectionObj.socket.SendBytes(msg.([]byte), 0)
	default:
		return false
	}

	return errors.IsError(err)
}

// Recv A function that synchronously receives a msg from the connection
func (connectionObj *Connection) Recv(flags zmq4.Flag) (interface{}, bool) {
	msg, err := connectionObj.socket.Recv(flags)

	return msg, errors.IsError(err)
}

// BuildConnectionString A function to build the connection string
func BuildConnectionString(ip string, port string) string {
	return fmt.Sprintf("tcp://%s:%s", ip, port)
}
