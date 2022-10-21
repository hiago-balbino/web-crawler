package main

import "github.com/hiago-balbino/web-crawler/api"

func main() {
	server := api.NewServer()
	server.Start()
}
