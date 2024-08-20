package main

import "github.com/TiagoAmaralFerreira/api-go/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
