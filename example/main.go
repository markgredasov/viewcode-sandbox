package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func main() {
	ctx := context.Background()

	sandboxConfig := sandbox_config.NewConfigMust()
	adapter, err := sandbox.NewAdapter(sandboxConfig)
	if err != nil {
		log.Fatal(err)
	}

	exec, err := model.NewExecution("go", `
	package main

	import "fmt"
	
	func main(){
		fmt.Println("4 + 5", "=", 9)
	}`)
	if err != nil {
		log.Fatal(err)
	}

	res, err := adapter.RunCode(ctx, exec)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println(res)
}
