package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

const _configPath = "config.yml"

// TestConfigPath is file path used for unit tests.
var TestConfigPath string

// AppConfig is the package level variable exposing the applications configs.
var AppConfig Config

// Config is a struct in yaml format for storing configs for the application.
type Config struct {
	Application struct {
		Name        string `default:"cs-code-review-bot" yaml:"name"`
		Environment string `default:"production" yaml:"environment"`
	} `yaml:"application"`
	Server struct {
		Port string `default:"8080" yaml:"port"`
	} `yaml:"server"`
	Github struct {
		GitHubURL   string `default:"https://api.github.com" yaml:"github_url"`
		GitHubToken string `default:"" yaml:"github_token"`
	} `yaml:"github"`
	Team []TeamSettings `yaml:"team_name"`
}

// TeamSettings holds the team settings
type TeamSettings struct {
	Name                string   `yaml:"name"`
	Channel             string   `yaml:"channel"`
	ChannelType         string   `yaml:"channel_type"`
	EnableCron          bool     `default:"false" yaml:"enable_cron"`
	CronElapsedDuration int      `default:"14400000000000" yaml:"cron_elapsed_duration"`
	Org                 string   `yaml:"org"`
	Repos               []string `yaml:"repos"`
	IgnoreUsers         []string `yaml:"ignore_users"`
}

// readFile opens the config file, parses it, and loads the Config struct
func readFile() {
	var file *os.File
	var err error
	if TestConfigPath == "" {
		file, err = os.Open(_configPath)
	} else {
		log.Debug("test config path is populated")
		file, err = os.Open(TestConfigPath)
	}
	if err != nil {
		if err.Error() == "open config.yml: no such file or directory" {
			log.Warn("Did not find config file. Proceeding with default config.")
			readEnv()
			return
		}
		processError(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Error closing file: %s\n", err)
		}
	}()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		processError(err)
	}
}

// readEnv reads local environment variables
func readEnv() {
	err := envconfig.Process("", &AppConfig)
	if err != nil {
		processError(err)
	}
}

// processError handles errors loading the config
func processError(err error) {
	log.Errorf("Could not load config. : %v", err)
	os.Exit(2)
}

// func setupLogging() {
// 	log.SetFormatter(&log.JSONFormatter{})
// 	log.SetLevel(log.DebugLevel)
// 	log.SetReportCaller(true)
// }

// LoadConfiguration is for loading config from files and environment and setting it ready to be read.
func LoadConfiguration() {
	// setupLogging()
	readFile()
	readEnv()
}
