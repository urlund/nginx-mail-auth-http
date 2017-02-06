# nginx-mail-auth-http
An HTTP authentication server to `ngx_mail_auth_http_module`. Auth lookups will be cached to minimize workload.

### command options

```
$ nginx-mail-auth-http -help
Usage of nginx-mail-auth-http:
  -auth-header string
    	Checks the specified header in requests sent to the authentication server (default "Auth-Key")
  -auth-key string
    	This header can be used as the shared secret to verify that the request comes from nginx
  -cache-cleanup string
    	Interval between cache cleanups (see: https://golang.org/pkg/time/#ParseDuration) (default "1m")
  -cache-ttl string
    	Time to keep proxy configs in cache since last usage (see: https://golang.org/pkg/time/#ParseDuration) (default "24h")
  -config-file string
    	Name of config file (default "config.json")
  -config-path string
    	Path where '-config-file' (and conf.d) can be found (default "/etc/nginx-mail-auth-http")
  -listen string
    	Address to handle requests on incoming connections (default ":8278")
```

## configuration
A configuration is based on (up to) 3 parts, *default*, *templates* and *proxy config* (only *default* is **required**)


### config.json

```
{
    "default": {
        PROXY_CONFIG
    },
    "templates": {
        "your-template-name": {
            PROXY_CONFIG
        },
        "your-template-name-N": {
            PROXY_CONFIG
        }
    }
}
```

#### proxy config

```
{
    "pop3": {
        "ip": "YOUR_POP3_IP",
        "port": YOUR_POP3_PORT
    },
    "imap": {
        "ip": "YOUR_IMAP_IP",
        "port": YOUR_IMAP_PORT
    },
    "smtp": {
        "ip": "YOUR_SMTP_IP",
        "port": YOUR_SMTP_PORT
    }
}
```

#### domain config
almost the same as **proxy config* but it supports *template*:

```
{
    "template": TEMPLATE_NAME,
    ...
}
```

If *template* is used *default* and *template* configuration can be overridden using *pop3*, *imap* and *smtp* settings as in proxy config (see configuration example 3).

## configuration examples

### example 1
A basic `config.json` example that will auth all domains to a single server.

```
{
    "defaults": {
        "pop3": {
            "ip": "YOUR_POP3_IP",
            "port": YOUR_POP3_PORT
        },
        "imap": {
            "ip": "YOUR_IMAP_IP",
            "port": YOUR_IMAP_PORT
        },
        "smtp": {
            "ip": "YOUR_SMTP_IP",
            "port": YOUR_SMTP_PORT
        }
    }
}
```

### example 2
Based on example 1 you can specify seperate domains to be auth'ed to another server by creating a seperate configuration in your `conf.d` folder (eg. `conf.d/example.com`):

```
{
    "pop3": {
        "ip": "ANOTHER_POP3_IP",
    },
    "imap": {
        "ip": "ANOTHER_IMAP_IP",
    },
    "smtp": {
        "ip": "ANOTHER_SMTP_IP",
    }
}
```

(Please note that this configuration does not contain any port configuration. They will be applied from "defaults" defined in `config.json`).

### example 3
If you have a lot of domains using the server from example 2, it might be a good idea to define the example 2 configuration as a `serverX` template in `config.json`

```
{
    "defaults": {
        ...
    },
    "templates": {
        "serverX": {
            ...
        }
    }
}
```

Instead of the complete configuration in your `conf.d/example.com` you can now define a reference to your template instead:

```
{
    "template": "serverX"
}
```

And even template references can be overridden in your domain configuration:

```
{
    "template": "serverX",
    "smtp": {
        "ip": "YET_ANOTHER_SMTP_IP"
    }
}
```

## NGINX configuration

`ngx_mail_auth_http_module` configuration:

```
mail {
    auth_http http://SERVER_IP:8278;
    
    # you are encuraged to configure auth_http_header (not required)
    # if you do so - remember to configure the '-auth-key' flag
    auth_http_header X-Auth-Key "YOUR_SECRET_STRING";
}
```

See [`ngx_mail_auth_http_module`](http://nginx.org/en/docs/mail/ngx_mail_auth_http_module.html) for more detailed configuration description.