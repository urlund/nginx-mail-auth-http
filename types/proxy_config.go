package types

import (
	"time"
)

type ProxyConfig struct {
	POP3     *ProtocolConfig `json:"pop3"`
	IMAP     *ProtocolConfig `json:"imap"`
	SMTP     *ProtocolConfig `json:"smtp"`
	Template string          `json:"template"`
	Timeout  time.Time
}

func (this *ProxyConfig) Port(protocol string) int {
	switch protocol {
	case "pop3":
		if this.POP3 != nil && this.POP3.Port != 0 {
			return this.POP3.Port
		}
	case "imap":
		if this.IMAP != nil && this.IMAP.Port != 0 {
			return this.IMAP.Port
		}
	case "smtp":
		if this.SMTP != nil && this.SMTP.Port != 0 {
			return this.SMTP.Port
		}
	}

	return 0
}

func (this *ProxyConfig) IP(protocol string) string {
	switch protocol {
	case "pop3":
		if this.POP3 != nil && this.POP3.IP != "" {
			return this.POP3.IP
		}
	case "imap":
		if this.IMAP != nil && this.IMAP.IP != "" {
			return this.IMAP.IP
		}
	case "smtp":
		if this.SMTP != nil && this.SMTP.IP != "" {
			return this.SMTP.IP
		}
	}

	return ""
}

func (this *ProxyConfig) Apply(config *ProxyConfig) {
	if config.POP3 != nil {
		if config.POP3.IP != "" {
			this.POP3.IP = config.POP3.IP
		}

		if config.POP3.Port > 0 {
			this.POP3.Port = config.POP3.Port
		}
	}

	if config.IMAP != nil {
		if config.IMAP.IP != "" {
			this.IMAP.IP = config.IMAP.IP
		}

		if config.IMAP.Port > 0 {
			this.IMAP.Port = config.IMAP.Port
		}
	}

	if config.SMTP != nil {
		if config.SMTP.IP != "" {
			this.SMTP.IP = config.SMTP.IP
		}

		if config.SMTP.Port > 0 {
			this.SMTP.Port = config.SMTP.Port
		}
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
