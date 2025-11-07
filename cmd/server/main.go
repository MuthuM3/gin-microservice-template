package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/MuthuM3/gin-microservice-template/internal/config"
)

var (
	version   = "0.0.1"
	buildTime = "Unknown"
	gitCommit = "Unknown"
)

func main() {
	// Parse command line flags
	var (
		configPath = flag.String("config", "", "path to configuration file")
		envPath    = flag.String("env", "development", "Environment (development|production)")
		showVer    = flag.Bool("version", false, "show version information")
	)

	flag.Parse()

	// show version information
	if *showVer {
		showVersion()
		return
	}

	cfg, err := LoadConfig(*configPath, *envPath)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := initLogger(cfg)
	logger.Printf("Starting Todo API %s in %s mode", version, *envPath)

	
}

func showVersion() {
	fmt.Printf("Todo API: %s\n", version)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// Load Configuration
func LoadConfig(configPath, env string) (*config.Config, error) {
	if configPath != "" {
		return config.Load(configPath)
	}
	return config.LoadForEnvironment(env)
}

func initLogger(cfg *config.Config) *log.Logger {
	flags := log.LstdFlags
	if cfg.Server.IsDevelopment() {
		flags |= log.Lshortfile
	}
	return log.New(os.Stdout, "[TODO-API]", flags)
}