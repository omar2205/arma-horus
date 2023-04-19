package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"oskr.nl/arma-horus.go/internal/jsonlog"
)

var (
	version   string
	buildTime string
)

type config struct {
	port int
	env  string // dev | staging | prod
	cors struct {
		trustedOrigins []string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	Server_folder   string   `json:"server_folder"`
	Server_script   string   `json:"server_script"`
	Server_pid_file string   `json:"server_pid_file"`
	Db_file         string   `json:"db_file"`
	Admins_email    []string `json:"admins_email"`
}

type application struct {
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func (c *config) load(logger *jsonlog.Logger) {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		logger.PrintFatal(err, map[string]string{
			"error": "error reading config file",
		})
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		logger.PrintFatal(err, map[string]string{
			"error": "error parsing config file",
		})
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trust CORS origins (space seperated)",
		func(urls string) error {
			cfg.cors.trustedOrigins = strings.Fields(urls)
			return nil
		},
	)

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	expvar.NewString("version").Set(version)
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	cfg.load(logger)
	logger.PrintInfo("Loaded config", map[string]string{
		"Server folder":   cfg.Server_folder,
		"Server script":   cfg.Server_script,
		"Server PID file": cfg.Server_pid_file,
	})

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
