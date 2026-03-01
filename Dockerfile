# renovate: datasource=docker depName=golang
FROM golang:1.20-alpine AS build

WORKDIR /workspace

COPY . .
ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o vultr-cloud-controller-manager .

# renovate: datasource=docker depName=alpine
FROM alpine:3.21.3
# renovate: datasource=repology depName=alpine_3_21/ca-certificates
RUN apk add --no-cache ca-certificates=20250911-r0

COPY --from=build /workspace/vultr-cloud-controller-manager /vultr-cloud-controller-manager
ENTRYPOINT ["/vultr-cloud-controller-manager"]
