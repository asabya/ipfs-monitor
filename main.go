package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sab94/ipfs-monitor/app"
)

func main() {

	// Bootstrap monitor app
	monitor, err := app.Bootstrap(context.Background())
	if err != nil {
		fmt.Printf("\n %v\n", err)
		os.Exit(1)
	}

	// Start tview app
	if err := monitor.App.Run(); err != nil {
		fmt.Printf("\n %v\n", err)
		os.Exit(1)
	}
}
