# http from scratch (go)

a tiny http server i built while learning go fundamentals.  

## setup

1. clone the repo and enter it.
2. (optional) download the demo video file if you do not already have it:

```bash
mkdir -p assets
curl -o assets/vim.mp4 https://storage.googleapis.com/qvault-webapp-dynamic-assets/lesson_videos/vim-vs-neovim-prime.mp4
```

## run the server

in terminal 1:

```bash
go run cmd/httpserver/main.go
```

server listens on: http://localhost:42069

## try requests

in terminal 2, run these one by one.

simple html responses:

```bash
curl -i http://localhost:42069/
curl -i http://localhost:42069/yourproblem
curl -i http://localhost:42069/myproblem
```

chunked streaming response (proxied from httpbin):

```bash
curl -i http://localhost:42069/httpbin/stream/20
```

video response:

```bash
curl -i http://localhost:42069/video --output vim.mp4
```

note: you can also open http://localhost:42069/video in your browser and the video should play directly.

## run tests

```bash
go test ./...
```

## reference

- [boot.dev lesson](https://www.boot.dev/lessons/b0cebf37-7151-48db-ad8a-0f9399f94c58)