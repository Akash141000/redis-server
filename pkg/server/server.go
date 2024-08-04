package server

import (
	"net"
	"redis/pkg/peer"
	"redis/pkg/store"

	"golang.org/x/exp/slog"
)

type Config struct {
	ListenAddr string
}

type ConfigOpts func(*Config)

type Server struct {
	*Config
	ln        net.Listener
	peers     map[*peer.Peer]bool
	addPeerch chan *peer.Peer
	quitch    chan struct{}
	store     store.Storer
}

func New(configOpts ...ConfigOpts) *Server {
	s := &Server{
		Config: &Config{
			ListenAddr: ":3000",
		},
		peers:     make(map[*peer.Peer]bool),
		quitch:    make(chan struct{}),
		addPeerch: make(chan *peer.Peer),
		store:     store.NewMemoryStore(),
	}

	for _, opt := range configOpts {
		opt(s.Config)
	}

	return s
}

func (s *Server) Start() error {
	slog.Info("server", "start", "server")

	// create channel to listen to error
	serverErr := make(chan error)

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go func() {
		// accept new peers
		if err := s.acceptPeersAndMessages(); err != nil {
			serverErr <- err
		}
	}()

	go func() {
		//
		if err := s.acceptConnections(); err != nil {
			serverErr <- err
		}
	}()

	//wait for the error and exit if error occurs
	err = <-serverErr
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) acceptConnections() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("server", "error accepting connection", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) acceptPeersAndMessages() error {
	for {
		select {
		case <-s.quitch:
			return nil
		case peer := <-s.addPeerch:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	//create new peer and add to server
	peer := peer.New(conn, s.store)
	s.addPeerch <- peer

	slog.Info("server", "new peer connected! remoteAddr", conn.RemoteAddr())

	//keep reading peer
	err := peer.ReadLoop()

	if err.Error() == "EOF" {
		slog.Info("server", "peer connection closed", err, "remoteAddr", peer.Conn.RemoteAddr())
		return nil
	}
	if err != nil {
		slog.Error("server", "peer readloop error", err, "remoteAddr", conn.RemoteAddr())
		return err
	}

	return nil
}

func WithListenAddr(listenAddr string) ConfigOpts {
	return func(cfg *Config) {
		cfg.ListenAddr = listenAddr
	}
}
