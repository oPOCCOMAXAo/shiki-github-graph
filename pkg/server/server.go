package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/opoccomaxao/shiki-github-graph/pkg/app"
	"github.com/opoccomaxao/shiki-github-graph/pkg/shiki"
	"github.com/opoccomaxao/shiki-github-graph/pkg/storage"
)

type Server struct {
	config Config

	Storage *storage.Storage
	Shiki   *shiki.Service

	listener           net.Listener
	usersProcessNotify chan struct{}
	animeProcessNotify chan struct{}
	warningLogger      *log.Logger
}

type Config struct {
	Storage storage.Config
	Shiki   shiki.Config

	Port int `env:"PORT" envDefault:"8080"`
}

func New(config Config) (*Server, error) {
	res := Server{
		config:             config,
		usersProcessNotify: make(chan struct{}, 100),
		animeProcessNotify: make(chan struct{}, 100),
	}

	var err error

	res.warningLogger, err = app.FileLogger("./warning.log")
	if err != nil {
		return nil, errors.WithMessage(err, "warningLogger init failed")
	}

	res.Storage, err = storage.New(config.Storage)
	if err != nil {
		return nil, errors.WithMessage(err, "Storage init failed")
	}

	res.Shiki, err = shiki.New(config.Shiki)
	if err != nil {
		return nil, errors.WithMessage(err, "Shiki init failed")
	}

	return &res, res.init()
}

func (s *Server) Serve(ctx context.Context) error {
	if s.listener != nil {
		return app.ErrAlreadyStarted
	}

	engine := gin.New()
	engine.SetTrustedProxies(nil)
	engine.GET("/user/:nick", s.GetUser)

	image := engine.Group("/image")
	image.Static("/", "data")

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.config.Port))
	if err != nil {
		return err
	}

	s.listener = listener

	go s.processUsersUpdate(ctx)
	go s.processAnimeUpdate(ctx)

	return http.Serve(listener, engine)
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) init() error {
	err := s.initImages()
	if err != nil {
		return err
	}

	return nil
}
