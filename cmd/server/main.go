package main

import (
	"flag"
	"fmt"
	"runtime"
)

var (
	version   = "0.0.1"
	buildTime = "Unknown"
	gitCommit = "Unknown"
)

func main() {
	// Parse command line flags
	var (
		// configPath = flag.String("config", "", "path to configuration file")
		// envPath    = flag.String("env", "development", "Environment (development|production)")
		showVer = flag.Bool("version", false, "show version information")
	)

	flag.Parse()

	// show version information
	if *showVer {
		showVersion()
		return
	}

}

func showVersion() {
	fmt.Printf("Todo API: %s\n", version)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
