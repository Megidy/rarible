package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Megidy/rarible/docs"

	"github.com/Megidy/rarible/cmd/app"
	"github.com/rs/zerolog/log"
)

// @title rarible client api
// @version 1.0
// @description swagger doucmentation.

// @host {base_url}
// @BasePath /v1
func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create application instance")
	}

	errsCh := app.Run()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errsCh:
		log.Err(err).Msg("Failed to run application")
	case <-sigs:
		log.Info().Msg("Received signal to shutdown")
	}

	err = app.Shutdown()
	if err != nil {
		log.Err(err).Msg("Failed to shutdown application")
	}

}
