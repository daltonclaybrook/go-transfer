
FROM alpine:latest

MAINTAINER Dalton Claybrook <daltonclaybrook@gmail.com>

WORKDIR "/opt"

ADD .docker_build/go-transfer /opt/bin/go-transfer

CMD ["/opt/bin/go-transfer"]