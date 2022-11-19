package main

import "github.com/hiago-balbino/web-crawler/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
