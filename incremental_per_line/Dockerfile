FROM alpine:3.7
RUN apk add --no-cache openssh ca-certificates

RUN mkdir /app
RUN mkdir /data
WORKDIR /app
COPY decrypt-incremental /app
COPY test.txt /app
ENTRYPOINT ["./decrypt-incremental"]
