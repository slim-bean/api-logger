package main

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/robfig/cron/v3"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	level.Info(logger).Log("msg", "starting")

	// Get the API URL from an environment variable
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		level.Error(logger).Log("msg", "API_URL environment variable is required")
		os.Exit(1)
	}
	cronString := os.Getenv("CRON")
	if cronString == "" {
		level.Error(logger).Log("msg", "CRON environment variable is required")
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "setting up cron", "api", apiURL, "cron", cronString)

	c := cron.New()
	_, err := c.AddFunc(cronString, func() {

		level.Info(logger).Log("msg", "querying API")

		response, err := http.Get(apiURL)
		if err != nil {
			level.Error(logger).Log("msg", "failed to query api", "err", err)
			return
		}
		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			level.Error(logger).Log("msg", "failed to read body", "err", err)
			return
		}

		// Log the response body
		level.Info(logger).Log("msg", "api response", "response", string(body))
	})

	if err != nil {
		level.Error(logger).Log("msg", "error creating cron job", "err", err)
		os.Exit(1)

	}

	c.Start()

	// catch signals and terminate the app
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sigc
	level.Info(logger).Log("msg", "exiting", "signal", s)
	os.Exit(0)

}
