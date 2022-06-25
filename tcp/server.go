package tcp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

const maxMsgLen = 5 * 1000 * 1000

// Conn returns messages parsed from TCP stream.
// This can be optimized by returning readers instead of byte array.
// But we keep bytes for simplicity.
type Conn interface {
	// ReadMessage is unsafe to call concurrently.
	ReadMessage() ([]byte, error)
	// WriteMessage is unsafe to call concurrently.
	WriteMessage([]byte) error
}

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}
func (c *Connection) ReadMessage() ([]byte, error) {
	sizeBuf := make([]byte, 4)
	_, err := io.ReadFull(c.conn, sizeBuf)
	if err != nil {
		return nil, err
	}
	size := binary.LittleEndian.Uint32(sizeBuf)
	if size > maxMsgLen {
		return nil, errors.New(fmt.Sprintf("Message size is more than %v", maxMsgLen))
	}
	message := make([]byte, size)
	_, err = io.ReadFull(c.conn, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Connection) WriteMessage(message []byte) error {
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(message)))
	_, err := c.conn.Write(size)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

type Handler interface {
	ServeTCP(Conn)
}

type Server struct {
	Handler Handler
	Addr    string
}

func ListenAndServe(ctx context.Context, addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe(ctx)
}

func (srv *Server) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		srv.Handler.ServeTCP(&Connection{conn: c})
		c.Close()
	}
}
