package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ghmeier/bloodlines/config"
	"github.com/jonnykry/expresso-billing/router"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	config, err := config.Init(path.Join(dir, "config.json"))
	if err != nil {
		fmt.Println("ERROR: unable to load config")
		return
	}

	b, err := router.New(config)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	fmt.Printf("Billing running on %s\n", config.Port)
	b.Start(":" + config.Port)
}
