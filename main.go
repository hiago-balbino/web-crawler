package main

import "github.com/hiago-balbino/web-crawler/internal/pkg/api"

func main() {
	server := api.NewServer()
	server.Start()
}
