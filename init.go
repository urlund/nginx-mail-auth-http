package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/urlund/nginx-mail-auth-http/types"
)

var (
	cacheTTL         string
	cacheCleanup     string
	configFile       string
	configPath       string
	listen           string
	authKey          string
	authHeader       string
	version          bool
	config           *types.Config
	proxyConfigCache map[string]*types.ProxyConfig
	timeout          time.Duration
	cleanup          time.Duration
)

func init() {
	proxyConfigCache = map[string]*types.ProxyConfig{}

	// why 8278? because is an unsigned port @ iana.org
	// (and, 25 + 110 + 143 = 278 and starting with 8 because 8080 is starting with 8 and typically used for a personally hosted web server)
	flag.StringVar(&listen, "listen", ":8278", "Address to handle requests on incoming connections")
	flag.StringVar(&cacheTTL, "cache-ttl", "24h", "Time to keep proxy configs in cache since last usage (see: https://golang.org/pkg/time/#ParseDuration)")
	flag.StringVar(&cacheCleanup, "cache-cleanup", "1m", "Interval between cache cleanups (see: https://golang.org/pkg/time/#ParseDuration)")
	flag.StringVar(&configFile, "config-file", "config.json", "Name of config file")
	flag.StringVar(&configPath, "config-path", "/etc/nginx-mail-auth-http", "Path where '-config-file' (and conf.d) can be found")
	flag.StringVar(&authKey, "auth-key", "", "This header can be used as the shared secret to verify that the request comes from nginx")
	flag.StringVar(&authHeader, "auth-header", "Auth-Key", "Checks the specified header in requests sent to the authentication server")
	flag.BoolVar(&version, "version", version, "Show version")
	flag.Parse()

	if version {
		fmt.Println("version: 1.0.0")
		os.Exit(0)
	}

	var err error
	timeout, err = time.ParseDuration(cacheTTL)
	if err != nil {
		panic(err)
	}

	cleanup, err = time.ParseDuration(cacheCleanup)
	if err != nil {
		panic(err)
	}

	jsonBlob, err := ioutil.ReadFile(filepath.Join(configPath, configFile))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonBlob, &config)
	if err != nil {
		panic(err)
	}

	if config.Default == nil {
		panic("default config is required")
	}
}
