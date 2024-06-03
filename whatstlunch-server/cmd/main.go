package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/ocxide/whatstlunch/cmd/cli"
	"github.com/ocxide/whatstlunch/cmd/config"
	"github.com/ocxide/whatstlunch/cmd/endpoints/dishes"
	"github.com/ocxide/whatstlunch/cmd/endpoints/infer"
)

// configPath - the path to the config file, can be empty
func readConfig(configPath string) (config.Config, error) {
	if configPath == "" {
		configPath = "config.toml"
	}

	raw, err := os.ReadFile(configPath)

	defaultConfig := config.Config{
		PublicDir: "public",
		Host:      "127.0.0.1:3456",
		Ai: config.AiConfig{
			Model: "llava:7b",
		},
	}

	if os.IsNotExist(err) {
		return defaultConfig, nil
	}

	if err != nil {
		return config.Config{}, err
	}

	cfg := config.Config{}
	_, err = toml.Decode(string(raw), &cfg)
	if err != nil {
		return config.Config{}, err
	}

	// Maybe should not check each field
	if cfg.PublicDir == "" {
		cfg.PublicDir = defaultConfig.PublicDir
	}

	if cfg.Host == "" {
		cfg.Host = defaultConfig.Host
	}

	if cfg.Ai.Model == "" {
		cfg.Ai.Model = defaultConfig.Ai.Model
	}

	return cfg, nil
}

func listen(host string, publicDir string) {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir(publicDir)))

	mux.HandleFunc("GET /dishes", dishes.Search)
	mux.HandleFunc("POST /infer-ingredients", infer.InferIngredients)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		mux.ServeHTTP(w, req)
	})

	fmt.Printf("Listening on: %s\n", host)

	server := http.Server{
		Addr:    host,
		Handler: handler,
	}

	err := server.ListenAndServe()
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
			configPath := cmd.Flags().Lookup("config").Value.String()

			config, err := readConfig(configPath)
			if err != nil {
				log.Fatal(err)
			}

			host := cmd.Flags().Lookup("host").Value.String()
			if host != "" {
				config.Host = host
			}

			publicDir := cmd.Flags().Lookup("public").Value.String()
			if publicDir != "" {
				config.PublicDir = publicDir
			}

			listen(config.Host, config.PublicDir)
		},
	}

	launchCmd.PersistentFlags().String("host", "", "Host to serve on")
	launchCmd.PersistentFlags().String("public", "", "Public directory")
	launchCmd.PersistentFlags().String("config", "", "Config file")

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
