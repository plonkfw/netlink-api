# Use build env
FROM golang:1.18.3-alpine3.16 as builder
LABEL maintainer="pers.edwin@gmail.com"

COPY ./src /src

RUN set -exu \
	&& mkdir /build \
	&& cd /src \
	&& CGO_ENABLED=0 \
			GOOS=linux \
			GOARCH=amd64 go build -a -v -o ../build .

# Switch to runtime env
FROM scratch

COPY --from=builder /build/netlink-api /netlink-api/netlink-api

ENV LISTEN="localhost:4821"
EXPOSE 4821

WORKDIR /netlink-api/
CMD ["./netlink-api"]
