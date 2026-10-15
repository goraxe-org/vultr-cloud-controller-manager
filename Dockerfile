# renovate: datasource=docker depName=golang
FROM golang:1.26-alpine AS build

WORKDIR /workspace

COPY . .
ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -trimpath -ldflags "-s -w -X main.version=$VERSION" -o vultr-cloud-controller-manager .

# renovate: datasource=docker depName=alpine
FROM alpine:3.21.8
# Pin OS packages to their current patched versions explicitly: the alpine:3.21.8
# base image is a frozen snapshot and can lag behind the v3.21 apk repo's latest
# security fixes (e.g. openssl CVE-2026-84782, fixed in 3.3.7-r2 but not yet baked
# into this tag).
# renovate: datasource=repology depName=alpine_3_21/ca-certificates
RUN apk add --no-cache ca-certificates=20260909-r0 \
# renovate: datasource=repology depName=alpine_3_21/openssl
	libssl3=3.3.7-r2 libcrypto3=3.3.7-r2 \
# renovate: datasource=repology depName=alpine_3_21/musl
	musl=1.2.5-r11 musl-utils=1.2.5-r11 \
# renovate: datasource=repology depName=alpine_3_21/zlib
	zlib=1.3.2-r1

COPY --from=build /workspace/vultr-cloud-controller-manager /vultr-cloud-controller-manager
ENTRYPOINT ["/vultr-cloud-controller-manager"]
