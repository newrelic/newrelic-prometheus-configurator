FROM alpine:3.16.0

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache --upgrade && apk add --no-cache tini curl bind-tools

COPY bin/configurator-${TARGETOS}-${TARGETARCH} /

RUN mv /configurator-${TARGETOS}-${TARGETARCH} /configurator && \
    chmod 755 /configurator
