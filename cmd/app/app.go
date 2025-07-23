package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Megidy/rarible/internal/client"
	"github.com/Megidy/rarible/internal/config"
	"github.com/Megidy/rarible/internal/handler"
	"github.com/Megidy/rarible/internal/service"
	httpserver "github.com/Megidy/rarible/pkg/servers/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const baseRaribleURL = "https://api.rarible.org/v0.1"

type App interface {
	Run() <-chan error
	Shutdown() error
}

type app struct {
	httpServer *httpserver.HttpServer
}

func NewApp() (App, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create new config: %w", err)
	}

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		log.Warn().Msgf("Invalid log level %s, defaulting to debug", cfg.LogLevel)
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	port := cfg.HttpServerPort
	if port != "" && port[0] != ':' {
		port = ":" + port
	}
	httpServer := httpserver.NewHttpServer(port)

	raribleClient := client.NewRaribleClient(cfg.RaribleApiKey, baseRaribleURL)

	nftService := service.NewNFTService(raribleClient)

	nftHandler := handler.NewNFTHandler(nftService)

	router := handler.NewRouter(httpServer.Echo, nftHandler)
	router.RegisterRoutes()

	return &app{
		httpServer: httpServer,
	}, nil
}

func (a *app) Run() <-chan error {
	errsCh := make(chan error, 1)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.httpServer.Run()
		if err != nil {
			errsCh <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errsCh)
	}()

	return errsCh

}
func (a *app) Shutdown() error {
	err := a.httpServer.Shutdown()
	if err != nil {
		return err
	}

	log.Info().Msgf("Successfully shut down HTTP server")
	return nil
}
