package lambda

import "github.com/kanaru0928/cms/internal/handler"

func init() {
	handler.Init()
}

func main() {
	handler.DefineListener(3000)
}
