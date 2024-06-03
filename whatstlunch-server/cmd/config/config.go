package config

type Config struct {
	PublicDir string   `toml:"public_dir"`
	Host      string   `toml:"host"`
	Ai        AiConfig `toml:"ai"`
}

type AiConfig struct {
	Model string `toml:"model"`
}
