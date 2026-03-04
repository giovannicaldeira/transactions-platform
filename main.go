package main

import "github.com/transactions-platform/cmd"

// @title           Transactions Platform API
// @version         1.0
// @description     A Go-based API platform built with Gin framework

// @host      localhost:8080
// @BasePath  /

// @schemes http https
func main() {
	cmd.Execute()
}
