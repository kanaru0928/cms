package main

import "github.com/kanaru0928/cms/internal/handler"

func main() {
	handler.Init()
	handler.DefineListener(3020)
}
