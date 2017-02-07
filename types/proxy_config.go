package types

import (
	"time"
)

type ProxyConfig struct {
	POP3     ProtocolConfig `json:"pop3,omitempty"`
	IMAP     ProtocolConfig `json:"imap,omitempty"`
	SMTP     ProtocolConfig `json:"smtp,omitempty"`
	Template string         `json:"template,omitempty"`
	Timeout  time.Time
}

func (this *ProxyConfig) Port(protocol string) int {
	switch protocol {
	case "pop3":
		if this.POP3.Port != 0 {
			return this.POP3.Port
		}
	case "imap":
		if this.IMAP.Port != 0 {
			return this.IMAP.Port
		}
	case "smtp":
		if this.SMTP.Port != 0 {
			return this.SMTP.Port
		}
	}

	return 0
}

func (this *ProxyConfig) IP(protocol string) string {
	switch protocol {
	case "pop3":
		if this.POP3.IP != "" {
			return this.POP3.IP
		}
	case "imap":
		if this.IMAP.IP != "" {
			return this.IMAP.IP
		}
	case "smtp":
		if this.SMTP.IP != "" {
			return this.SMTP.IP
		}
	}

	return ""
}

func (this *ProxyConfig) Apply(config *ProxyConfig) {
	if config.POP3.IP != "" {
		this.POP3.IP = config.POP3.IP
	}

	if config.POP3.Port > 0 {
		this.POP3.Port = config.POP3.Port
	}

	if config.IMAP.IP != "" {
		this.IMAP.IP = config.IMAP.IP
	}

	if config.IMAP.Port > 0 {
		this.IMAP.Port = config.IMAP.Port
	}

	if config.SMTP.IP != "" {
		this.SMTP.IP = config.SMTP.IP
	}

	if config.SMTP.Port > 0 {
		this.SMTP.Port = config.SMTP.Port
	}

	if config.Template != "" {
		if config.Template != "" {
			this.Template = config.Template
		}

		if config.Template != "" {
			this.Template = config.Template
		}
	}
}
