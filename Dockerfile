FROM alpine as alpine
RUN apk add -U --no-cache ca-certificates

# builder
FROM scratch
ENTRYPOINT ["/gotify-bridge"]
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY gotify-bridge /