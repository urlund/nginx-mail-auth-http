package types

type Config struct {
	Default   *ProxyConfig            `json:"default"`
	Templates map[string]*ProxyConfig `json:"templates"`
	Ttl       string                  `json:"ttl"`
}
