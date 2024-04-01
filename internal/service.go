package internal

import (
	"fmt"
	"github.com/anycable/anycable-go/cli"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

type Service struct {
	config *Config
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) Run() int {
	handler := s.createHandler()
	server := NewServer(s.config, handler)
	upstream := NewUpstreamProcess(s.config.UpstreamCommand, s.config.UpstreamArgs...)

	server.Start()
	defer server.Stop()

	s.setEnvironment()

	exitCode, err := upstream.Run()
	if err != nil {
		panic(err)
	}

	return exitCode
}

// Private

func (s *Service) createHandler() http.Handler {
	targetUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%d", s.config.TargetPort))
	cache := NewMemoryCache(s.config.CacheSizeBytes, s.config.MaxCacheItemSizeBytes)
	options := HandlerOptions{
		cache:                    cache,
		targetUrl:                targetUrl,
		xSendfileEnabled:         s.config.XSendfileEnabled,
		maxCacheableResponseBody: s.config.MaxCacheItemSizeBytes,
		maxRequestBody:           s.config.MaxRequestBody,
		badGatewayPage:           s.config.BadGatewayPage,
	}

	handler := NewHandler(options)

	if !s.config.AnyCableDisabled {
		anycable, err := s.runAnyCable(slog.Default())
		if err != nil {
			panic(err)
		}
		handler = NewAnyCableHandler(anycable, handler)
	}

	return handler
}

func (s *Service) setEnvironment() {
	// Set PORT to be inherited by the upstream process.
	os.Setenv("PORT", fmt.Sprintf("%d", s.config.TargetPort))
}

func (s *Service) runAnyCable(l *slog.Logger) (*cli.Embedded, error) {
	c := cli.NewConfig()

	opts := []cli.Option{
		cli.WithName("AnyCable"),
		cli.WithDefaultRPCController(),
		cli.WithDefaultBroker(),
		cli.WithDefaultSubscriber(),
		cli.WithDefaultBroadcaster(),
		cli.WithLogger(l),
	}

	runner, err := cli.NewRunner(c, opts)

	if err != nil {
		return nil, err
	}

	anycable, err := runner.Embed()

	return anycable, err
}
