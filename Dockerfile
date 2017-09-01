FROM alpine:latest
MAINTAINER Gennaro Vietri <gennaro.vietri@bitbull.it>

RUN apk --update add socat

ADD main /
ADD index.html /

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 8888

CMD ["/main"]
