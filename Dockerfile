FROM alpine:3.20.1

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache --upgrade && apk add --no-cache tini curl bind-tools

COPY bin/prometheus-configurator-${TARGETOS}-${TARGETARCH} /

COPY THIRD_PARTY_NOTICES.md /

RUN mv /prometheus-configurator-${TARGETOS}-${TARGETARCH} /prometheus-configurator && \
    chmod 755 /prometheus-configurator

# creating the nri-agent user used only in unprivileged mode
RUN addgroup -g 2000 nri-agent && adduser -D -u 1000 -G nri-agent nri-agent

USER nri-agent

ENTRYPOINT ["/prometheus-configurator"]
