package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sab94/ipfs-monitor/app"
)

func main() {
	monitor, err := app.Start(context.Background())
	if err != nil {
		fmt.Printf("\n %v\n", err)
		os.Exit(1)
	}
	if err := monitor.App.Run(); err != nil {
		fmt.Printf("\n %v\n", err)
		os.Exit(1)
	}
}
