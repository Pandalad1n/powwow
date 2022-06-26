package tcp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
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

func Dial(address string) (Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &conn{conn: c}, nil
}

type Handler interface {
	ServeTCP(context.Context, Conn)
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
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		<-ctx.Done()
		_ = listener.Close()
		wg.Done()
	}()
	for {
		// Usually we want to limit how many connections can be served in parallel.
		// But we keep it simple.
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			srv.Handler.ServeTCP(ctx, &conn{conn: c})
			cancel()
			_ = c.Close()
			wg.Done()
		}()
	}
}

// conn converts endless stream into desecrate messages.
// It prepends each message with its length.
type conn struct {
	conn net.Conn
}

func (c *conn) ReadMessage() ([]byte, error) {
	sizeBuf := make([]byte, 4)
	_, err := io.ReadFull(c.conn, sizeBuf)
	if err != nil {
		return nil, err
	}
	size := binary.LittleEndian.Uint32(sizeBuf)
	if size > maxMsgLen {
		return nil, fmt.Errorf("message size is more than %d", maxMsgLen)
	}
	message := make([]byte, size)
	_, err = io.ReadFull(c.conn, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *conn) WriteMessage(message []byte) error {
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

func (c *conn) Close() error {
	return c.conn.Close()
}
