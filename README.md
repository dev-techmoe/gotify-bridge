# Gotify Bridge

[![goreleaser](https://github.com/dev-techmoe/gotify-bridge/actions/workflows/goreleaser.yml/badge.svg?branch=master)](https://github.com/dev-techmoe/gotify-bridge/actions/workflows/goreleaser.yml)
[![GitHub license](https://img.shields.io/github/license/dev-techmoe/gotify-bridge)](https://github.com/dev-techmoe/gotify-bridge/blob/master/LICENSE)

A bridge can helps you devlier your messages on Gotify to other place.  
Current only support WebPush

## What is Gotify

Gotify is a message server written in Golang. see [Homepage of Gotify](https://gotify.net/) for more detail.

## Download

### Binary

See [Release page](https://github.com/dev-techmoe/gotify-bridge/releases/tag/v1.0)

### Docker

```plain
docker pull ghcr.io/dev-techmoe/gotify-bridge
```

## Usage

1. write `config.json`, specify your gotify server address (websocket address) and webserver address (for showing a webpage to register webpush in browser)

   ```json
   {
     "Http": {
       "ListenAddress": "127.0.0.1:8080"
     },
     "Gotify": {
       "Address": "ws://your-gotify-address.com/stream?token=your_token_here"
     }
   }
   ```

   NOTICE: if your Gotify instance is run with TLS, use `wss://` instead of `ws://`

2. run Gotify bridge  
   if you use binary:
   ```bash
   gotify-bridge --config /path/to/your/config.json
   ```
   if you use docker (change the config path and port manually):
   ```bash
   docker run --name gotify-bridge \
        -v /path/to/your/config.json:/config.json \
        -p 8080:8080 \
        ghcr.io/dev-techmoe/gotify-bridge
   ```
3. visit `127.0.0.1:8080` (or other address you had set for webserver), click the "subscribe" button for register WebPush
4. Now try to push a message to your Gotify server for testing the bridge work normally.

## LICENSE

MIT
