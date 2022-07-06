FROM alpine:3.16.0

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache --upgrade && apk add --no-cache tini curl bind-tools

COPY bin/configurator-${TARGETOS}-${TARGETARCH} /
COPY configurator/testdata/input-test.yaml /

RUN mv /configurator-${TARGETOS}-${TARGETARCH} /configurator && \
    chmod 755 /configurator

ENTRYPOINT ["/sbin/tini", "--", "/configurator", "-input", "/input-test.yaml", "-output", "/etc/prometheus/config/config.yaml"]
