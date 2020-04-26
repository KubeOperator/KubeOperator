package main

import (
	"ko3-gin/internal/server"
)

func main() {
	if err := server.Start(); err != nil {
			panic("")
		}
}
