# STAGE: mibs
FROM docker.io/ciscocx/mibs:0.4.0 as mibs

# STAGE: build
FROM golang:1.13.1-stretch AS build
ADD . /src
WORKDIR /src

## Build of.
RUN make build

# STAGE: final
FROM debian:buster-slim

## Install mibs.
ENV SNMP_MIBS_DIR /mibs/mibs.snmplabs.com/json
WORKDIR $SNMP_MIBS_DIR
COPY --from=mibs $SNMP_MIBS_DIR .

## Install distro tools.
RUN apt-get update && apt-get --no-install-recommends -y install \
    # Add core programs \
    ca-certificates coreutils git make jq \
    # Add Debugging Tools (Removable in the future) \
    iproute2 net-tools nmap curl wget dnsutils && \
    apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

## Install tini.
RUN cd /tmp && \
    wget -O tini https://github.com/krallin/tini/releases/download/v0.18.0/tini-static-amd64 && \
    chmod 755 tini && \
    mv tini /sbin/tini

## Create a dedicated app user.
WORKDIR /app
RUN addgroup --gid 1000 --system app && \
    adduser  --uid 1000 --system --ingroup app --home /app --no-create-home app
USER app

## Install of.
COPY --from=build /src/of .

## Install a copy of the source code (re: licensing).
ADD . /src

## Configure runtime behavior.
ENV PATH="/app:${PATH}"
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["of"]
