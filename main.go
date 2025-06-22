package main

import (
	_ "github.com/Henrique-Batista/caddy-aspnetcore-adapter/adapter"
	_ "github.com/Henrique-Batista/caddy-aspnetcore-adapter/middleware"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
)

func main() {
	caddycmd.Main()
}
