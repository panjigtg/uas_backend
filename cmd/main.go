package main

import "uas/config"

func main() {
	app := config.Bootstrap()
	
	app.Listen(":3000")
}