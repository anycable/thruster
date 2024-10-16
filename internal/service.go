package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/anycable/anycable-go/cli"
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
	handlerOptions := HandlerOptions{
		cache:                    s.cache(),
		targetUrl:                s.targetUrl(),
		xSendfileEnabled:         s.config.XSendfileEnabled,
		maxCacheableResponseBody: s.config.MaxCacheItemSizeBytes,
		maxRequestBody:           s.config.MaxRequestBody,
		badGatewayPage:           s.config.BadGatewayPage,
		forwardHeaders:           s.config.ForwardHeaders,
	}

	handler := NewHandler(handlerOptions)
	handler = s.maybeHandleAnyCable(handler)
	server := NewServer(s.config, handler)
	upstream := NewUpstreamProcess(s.config.UpstreamCommand, s.config.UpstreamArgs...)

	server.Start()
	defer server.Stop()

	s.setEnvironment()

	exitCode, err := upstream.Run()
	if err != nil {
		slog.Error("Failed to start wrapped process", "command", s.config.UpstreamCommand, "args", s.config.UpstreamArgs, "error", err)
		return 1
	}

	return exitCode
}

// Private

func (s *Service) cache() Cache {
	return NewMemoryCache(s.config.CacheSizeBytes, s.config.MaxCacheItemSizeBytes)
}

func (s *Service) targetUrl() *url.URL {
	url, _ := url.Parse(fmt.Sprintf("http://localhost:%d", s.config.TargetPort))
	return url
}

func (s *Service) setEnvironment() {
	// Set PORT to be inherited by the upstream process.
	os.Setenv("PORT", fmt.Sprintf("%d", s.config.TargetPort))
}

func (s *Service) maybeHandleAnyCable(handler http.Handler) http.Handler {
	if !s.config.AnyCableDisabled {
		anycable, err := s.runAnyCable(slog.Default())
		if err != nil {
			panic(err)
		}
		handler = NewAnyCableHandler(anycable, handler)
	}

	return handler
}

func (s *Service) runAnyCable(l *slog.Logger) (*cli.Embedded, error) {
	argsWithProg := append([]string{"anycable-go"}, strings.Fields(s.config.AnyCableOptions)...)

	c, err, _ := cli.NewConfigFromCLI(argsWithProg)
	if err != nil {
		return nil, err
	}

	opts := []cli.Option{
		cli.WithName("AnyCable"),
		cli.WithDefaultRPCController(),
		cli.WithDefaultBroker(),
		cli.WithDefaultSubscriber(),
		cli.WithDefaultBroadcaster(),
		cli.WithLogger(l),
		cli.WithTelemetry("variant", "thruster"),
	}

	runner, err := cli.NewRunner(c, opts)

	if err != nil {
		return nil, err
	}

	anycable, err := runner.Embed()

	return anycable, err
}
