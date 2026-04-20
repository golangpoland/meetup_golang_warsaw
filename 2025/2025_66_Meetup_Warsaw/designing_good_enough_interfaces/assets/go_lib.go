package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var name string
	var cmd = &cobra.Command{
		Use: "greet",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s!\n", name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "World", "Name to greet")
	cmd.Execute()
}
