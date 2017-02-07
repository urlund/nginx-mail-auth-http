FROM debian:jessie

ADD https://github.com/urlund/nginx-mail-auth-http/releases/download/1.0.0/nginx-mail-auth-http

RUN chmod +x /usr/local/bin/nginx-mail-auth-http \
    mkdir -p /etc/nginx-mail-auth-http/conf.d/ \
    echo '{"default":{}}' > /etc/nginx-mail-auth-http/config.json

EXPOSE 8278

CMD ["nginx-mail-auth-http"]
