FROM golang:1.20-alpine AS build

WORKDIR /workspace

COPY . .
ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o vultr-cloud-controller-manager .

FROM alpine:3.19.0
RUN apk add --no-cache ca-certificates=20230506-r0

COPY --from=build /workspace/vultr-cloud-controller-manager /usr/local/bin/vultr-cloud-controller-manager
ENTRYPOINT ["vultr-cloud-controller-manager"]
