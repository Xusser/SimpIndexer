package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Host          string `json:"host"`
	Port          int64  `json:"port"`
	BaseURL       string `json:"base_url"`
	PathToExplore string `json:"explore"`
}

func NewDefaultConfig() *Config {
	result := &Config{
		Host:    "127.0.0.1",
		Port:    80,
		BaseURL: "",
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	result.PathToExplore = filepath.Dir(ex)

	return result
}

var _config *Config = NewDefaultConfig()

var (
	FlagDebugLog = false
	FlagTraceLog = false

	FlagHost          string = "0.0.0.0"
	FlagPort          int64  = 8000
	FlagPathToExplore string = ""
)

func init() {
	flag.BoolVar(&FlagDebugLog, "v", false, "Logging in debug mode")
	flag.BoolVar(&FlagTraceLog, "vv", false, "Logging in verbose mode")

	flag.StringVar(&FlagHost, "b", "0.0.0.0", "Host to bind")
	flag.Int64Var(&FlagPort, "p", 8000, "TCP port to bind")
	flag.StringVar(&FlagPathToExplore, "e", "", "Path to be explored, working directory will be used if it's unset")
}

func Init() error {

	if ip := net.ParseIP(FlagHost); ip == nil {
		return fmt.Errorf("invalid host")
	}
	_config.Host = FlagHost

	if FlagPort < 0 || FlagPort > 65535 {
		return fmt.Errorf("invalid port")
	}
	_config.Port = FlagPort

	if FlagPathToExplore != "" || flag.Arg(0) != "" {

		if flag.Arg(0) != "" {
			log.Debug().Msgf("Detected arg0: %s", flag.Arg(0))
			FlagPathToExplore = flag.Arg(0)
		}

		fi, err := os.Stat(FlagPathToExplore)
		if err != nil {
			return fmt.Errorf("invalid path to explore; %v", err)
		}

		if !fi.IsDir() {
			return fmt.Errorf("path should not be file")
		}

		_config.PathToExplore = FlagPathToExplore
	}

	return nil
}

func Get() *Config {
	return _config
}
