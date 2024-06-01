package main

import (
	"fmt"
	"net/http"

	"github.com/ocxide/whatstlunch/cmd/cli"
	"github.com/ocxide/whatstlunch/cmd/endpoints/dishes"
	"github.com/spf13/cobra"
)

func listen(host string) {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/dishes", dishes.Search)

	fmt.Printf("Listening on: %s\n", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "tlunch",
		Short: "Root command",
	}

	launchCmd := &cobra.Command{
		Use:   "launch",
		Short: "Launch command",
		Run: func(cmd *cobra.Command, args []string) {
			host := cmd.Flags().Lookup("host").Value.String()
			if host == "" {
				host = "127.0.0.1:3456"
			}

			listen(host)
		},
	}

	launchCmd.PersistentFlags().String("host", "", "Host to serve on")
	rootCmd.AddCommand(launchCmd)

	loadCmd := &cobra.Command{
		Use:   "load",
		Short: "Load command",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Load(args[0])
		},
	}

	loadCmd.Args = cobra.MinimumNArgs(1)
	rootCmd.AddCommand(loadCmd)

	rootCmd.Execute()
}
