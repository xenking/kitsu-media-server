package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a application configuration structure
type Config struct {
	Database struct {
		Host        string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
		Port        string `yaml:"port" env:"DB_PORT" env-description:"Database port"`
		Username    string `yaml:"username" env:"DB_USER" env-description:"Database user name"`
		Password    string `yaml:"password" env:"DB_PASSWORD" env-description:"Database user password"`
		Name        string `yaml:"db_name" env:"DB_NAME" env-description:"Database name"`
		Connections int    `yaml:"connections" env:"DB_CONNECTIONS" env-description:"Total number of database connections"`
	} `yaml:"database"`
	Default struct {
		UserImg string `yaml:"user_img" env-description:"Default user image"`
	} `yaml:"default"`
	Server struct {
		Host      string `yaml:"host" env:"SRV_HOST,HOST" env-description:"Server host" env-default:"localhost"`
		Port      string `yaml:"port" env:"SRV_PORT,PORT" env-description:"Server port" env-default:"8080"`
		JWTSecret string `yaml:"secret" env:"SRV_SECRET,SECRET" env-description:"JWT secret string"`
	} `yaml:"server"`
}

// args command-line parameters
type args struct {
	ConfigPath      string
	SessionProvider string
}

// processArgs processes and handles CLI arguments
func processArgs(cfg interface{}) args {
	var a args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "/config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}

var Global = &struct {
	JWTSecret []byte
	UserImg   string
}{}

func (cfg *Config) Init() {
	args := processArgs(cfg)
	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, cfg); err != nil {
		log.Println(err)
		os.Exit(2)
	}
	setGlobal(cfg)
}

func setGlobal(cfg *Config) {
	Global.JWTSecret = []byte(cfg.Server.JWTSecret)
	Global.UserImg = cfg.Default.UserImg
}
