package main

import (
	"zerologix-homework/cmd"
)

// @title zerologix-homework
// @version 1.0
// @description homework
// @BasePath /api/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
