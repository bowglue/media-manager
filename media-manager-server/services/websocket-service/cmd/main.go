package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
	mu    sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New client connected: ", ws.RemoteAddr())
	s.mu.Lock()
	s.conns[ws] = true
	s.mu.Unlock()

	// Clean up on disconnect
	defer func() {
		s.mu.Lock()
		delete(s.conns, ws)
		s.mu.Unlock()
		ws.Close()
		fmt.Println("Client disconnected:", ws.RemoteAddr())
	}()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(msg []byte) {
	for ws := range s.conns {
		go func() {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("write error:", err)
			}
		}()
	}
}

func main() {
	server := NewServer()
	log.Println("WebSocket server started on :15000")
	webSocketHandler := websocket.Server{
		Handler:   server.handleWS,
		Handshake: nil,
	}

	http.Handle("/ws", webSocketHandler)
	http.ListenAndServe(":15000", nil)

}
