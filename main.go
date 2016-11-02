package main

import (
	"fmt"
	"os"

	"github.com/jonnykry/expresso-billing/router"
)

func main() {
	b, err := router.New()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	port := os.Getenv("PORT");
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Billing running on %s\n", port)
	b.Start(":"+port)
}
