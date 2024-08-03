package server

import (
	"net"
	"redis/pkg/peer"
	"redis/pkg/proto"
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
	msgch     chan peer.Message
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
		msgch:     make(chan peer.Message),
		store:     store.NewMemoryStore(),
	}

	for _, opt := range configOpts {
		opt(s.Config)
	}

	return s
}

func (s *Server) Start() error {
	slog.Info("server", "start", "server")

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("server", "error starting the server", err)
		return err
	}

	s.ln = ln

	// accept new peers
	if err := s.acceptPeersAndMessages(); err != nil {
		return err
	}

	//
	if err := s.acceptConnections(); err != nil {
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
		case msg := <-s.msgch:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("server", "error handle message", err)
			}
		case peer := <-s.addPeerch:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	//create new peer and add to server
	peer := peer.New(conn, s.msgch)
	s.addPeerch <- peer

	slog.Info("server", "new peer connected!", "remoteAddr", conn.RemoteAddr())

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

func (s *Server) handleMessage(msg peer.Message) error {
	//parse the incoming message
	cmd, err := proto.ParseCommand(msg.Data)
	if err != nil {
		return err
	}

	// check if command is 'Set' or 'Get'
	switch v := cmd.(type) {
	case proto.SetCommand:
		return s.store.Set(string(v.Key), v.Value)
	case proto.GetCommand:
		val, err := s.store.Get(v.Key)
		if err != nil {
			return err
		}
		//send the value found for the key over connection
		msg.Peer.Send(val)
	}
	return nil
}

func WithListenAddr(listenAddr string) ConfigOpts {
	return func(cfg *Config) {
		cfg.ListenAddr = listenAddr
	}
}
