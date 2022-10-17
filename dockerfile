FROM golang:alpine AS build-env
RUN apk --no-cache add ca-certificates git
WORKDIR /go/src/rp
COPY . .
RUN CGO_ENABLED=0 go build -o output/rp

FROM debian:stable-slim
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /go/src/rp/output/rp /rp

ENV HTTP 0.0.0.0:80
ENV HTTPS 0.0.0.0:443
ENV SOURCE http://127.0.0.1:8080
ENV CRT ""
ENV KEY ""

ENTRYPOINT [ "bash", "-c", "/rp -http=${HTTP} -https=${HTTPS} -source=${SOURCE} -crt=\"${CRT}\" -key=\"${KEY}\"" ]