package main

import (
	"github.com/kanaru0928/cms/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	handler.Init()
	handler.DefineListener(3020)
}
