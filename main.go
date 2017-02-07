package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/urlund/nginx-mail-auth-http/types"
)

func debugConfig(p interface{}) {
	x, _ := json.Marshal(p)
	fmt.Println(string(x))
}

func handleResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Add("Auth-Status", message)

	if message == "OK" {
		if r.Header.Get("Auth-Method") == "cram-md5" {
			w.Header().Add("Auth-Pass", "plain-text-pass")
		}
	}
}

func getProxyConfig(domain string) (types.ProxyConfig, string) {
	proxyConfig := config.Default
	var domainConfig types.ProxyConfig

	jsonBlob, err := ioutil.ReadFile(filepath.Join(configPath, "conf.d", domain))
	if err != nil {
		return proxyConfig, ""
	}

	err = json.Unmarshal(jsonBlob, &domainConfig)
	if err != nil {
		return proxyConfig, "unable to load proxy config"
	}

	// check if we need to apply a template
	if domainConfig.Template != "" {
		if templateConfig, templateFound := config.Templates[domainConfig.Template]; templateFound == true {
			proxyConfig.Apply(&templateConfig)
		}
	}

	// ...
	proxyConfig.Apply(&domainConfig)

	return proxyConfig, ""
}

func getAuthServerAndPort(w http.ResponseWriter, r *http.Request, domain string) (err string) {
	proxyConfig, cacheFound := proxyConfigCache[domain]

	if cacheFound == false || (cacheFound == true && time.Now().After(proxyConfig.Timeout)) {
		// get proxy config
		proxyConfig, err = getProxyConfig(domain)
		if err != "" {
			return err
		}

		proxyConfigCache[domain] = proxyConfig
		w.Header().Add("X-Cache", "MISS")
	} else {
		w.Header().Add("X-Cache", "HIT")
	}

	// get auth ip and port
	protocol := r.Header.Get("Auth-Protocol")
	ip := proxyConfig.IP(protocol)
	port := proxyConfig.Port(protocol)

	// check if ip and port was found
	if ip == "" || port == 0 {
		return fmt.Sprintf("unable to find proxy server or port for protocol: '%s'", protocol)
	}

	// extend cache timeout
	proxyConfig.Timeout = time.Now().Add(timeout)
	proxyConfigCache[domain] = proxyConfig

	// set auth headers
	w.Header().Add("Auth-Server", ip)
	w.Header().Add("Auth-Port", strconv.Itoa(port))

	// no error occured
	return ""
}

func main() {
	// cleanup expired cache entries
	go func() {
		for {
			for domain, proxyConfig := range proxyConfigCache {
				if time.Now().After(proxyConfig.Timeout) {
					delete(proxyConfigCache, domain)
				}
			}

			time.Sleep(cleanup)
		}
	}()

	// handle mail proxy auth
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if authKey != "" && r.Header.Get(authHeader) != authKey {
			handleResponse(w, r, "invalid auth key, check your configuration")
			return
		}

		user := r.Header.Get("Auth-User")
		if user == "" || r.Header.Get("Auth-Pass") == "" {
			handleResponse(w, r, "username and password are required")
			return
		}

		// validate user as email address
		re := regexp.MustCompile("(.+)@(.+\\..+)")
		if re.Match([]byte(user)) == false {
			handleResponse(w, r, "please use a valid email address")
			return
		}

		// user parts consist of [email, name, domain]
		userParts := re.FindStringSubmatch(user)
		if len(userParts) != 3 {
			handleResponse(w, r, "invalid email address")
			return
		}

		// get domain proxy server and port
		err := getAuthServerAndPort(w, r, userParts[2])
		if err != "" {
			handleResponse(w, r, err)
			return
		}

		// ...
		handleResponse(w, r, "OK")
	})

	// start (keep things running)
	http.ListenAndServe(listen, nil)
}
