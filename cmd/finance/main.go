package main

import "github.com/TiagoAmaralFerreira/api-go/configs"

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	println(config.DBDriver)
}
