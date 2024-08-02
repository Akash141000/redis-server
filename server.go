package main

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/exp/slog"
)

type Config struct {
	ListenAddr string
}

type ConfigOpts func(*Config)

type Server struct {
	*Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerch chan *Peer
	quitch    chan struct{}
	msgch     chan []byte
	store     *MemoryStore
}

func NewServer(configOpts ...ConfigOpts) *Server {
	s := &Server{
		Config: &Config{
			ListenAddr: ":3000",
		},
		peers:     make(map[*Peer]bool),
		quitch:    make(chan struct{}),
		addPeerch: make(chan *Peer),
		msgch:     make(chan []byte),
		store:     NewMemoryStore(),
	}

	for _, opt := range configOpts {
		opt(s.Config)
	}

	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("server", "error starting the server", err)
		return err
	}

	s.ln = ln

	//
	go s.acceptNewPeers()

	slog.Info("server", "start", "server")

	return s.acceptConnection()
}

func (s *Server) acceptConnection() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("server", "error accepting connection", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) acceptNewPeers() error {
	for {
		select {
		case <-s.quitch:
			return nil
		case rawMsg := <-s.msgch:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("server", "error handle raw message", err)
			}
			fmt.Println("raw msg", rawMsg)
		case peer := <-s.addPeerch:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.msgch)
	s.addPeerch <- peer

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	err := peer.readLoop()
	if err == io.EOF {
		slog.Info("server", "peer connection closed", err, "remoteAddr", peer.conn.RemoteAddr())
		return
	}
	if err != nil {
		slog.Error("server", "peer readloop error", err, "remoteAddr", conn.RemoteAddr())
	}
}

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := ParseCommand(rawMsg)
	if err != nil {
		return err
	}
	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("trying to set command ", v)
	}
	return nil
}

func WithListenAddr(listenAddr string) ConfigOpts {
	return func(cfg *Config) {
		cfg.ListenAddr = listenAddr
	}
}
