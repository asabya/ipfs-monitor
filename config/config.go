package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	DefaultDir = "ipfs-monitor"
	ConfigFile = "config.yml"
)

// Common Module configs
type Common struct {
	PositionSettings
	Bordered        bool
	Enabled         bool
	RefreshInterval int
	Title           string
}

// PositionSettings places a module in a grid
type PositionSettings struct {
	Top    int `yaml:"top"`
	Left   int `yaml:"left"`
	Height int `yaml:"height"`
	Width  int `yaml:"width"`
}

// Settings is module settings
type Settings struct {
	Common *Common
	URL    string
}

// Config is monitor generated yml config
type Config struct {
	Monitor struct {
		Colors struct {
			Border struct {
				Focusable string `yaml:"focusable"`
				Focused   string `yaml:"focused"`
				Normal    string `yaml:"normal"`
			} `yaml:"border"`
			Background string `yaml:"background"`
			Text       string `yaml:"text"`
		} `yaml:"colors"`
		Grid struct {
			Columns    []int  `yaml:"columns"`
			Rows       []int  `yaml:"rows"`
			Background string `yaml:"background"`
			Border     bool   `yaml:"border"`
		} `yaml:"grid"`
		RefreshInterval int                      `yaml:"refreshInterval"`
		Widgets         map[string]WidgetConfigs `yaml:"widgets"`
	} `yaml:"monitor"`
}

// WidgetConfigs is config for each widgets
type WidgetConfigs struct {
	Enabled          bool             `yaml:"enabled"`
	PositionSettings PositionSettings `yaml:"position"`
	RefreshInterval  int              `yaml:"refreshInterval"`
	Title            string           `yaml:"title"`
}

// CreateOrLoadConfigFile creates or loads config file
// from config directory
func CreateOrLoadConfigFile() *Config {
	createConfigDir()
	filePath, err := CreateFile(ConfigFile)
	if err != nil {
		fmt.Println("Unable to create config file", err.Error())
		os.Exit(1)
	}

	// If the file is empty, write to it
	file, _ := os.Stat(filePath)

	if file.Size() == 0 {
		if ioutil.WriteFile(filePath, []byte(defaultConfigFile), 0600) != nil {
			fmt.Println("Unable to write to config file", err.Error())
			os.Exit(1)
		}
	}

	ymlCfg, err := parseYaml(filePath)
	if err != nil {
		fmt.Println("Unable to load config file", err.Error())
		os.Exit(1)
	}

	return ymlCfg
}

// CreateFile creates config file
func CreateFile(name string) (string, error) {
	appDir, err := configDir()
	if err != nil {
		return "", err
	}
	confFile := filepath.Join(appDir, name)

	_, err = os.Stat(confFile)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(confFile)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return confFile, nil
}

func createConfigDir() {
	configDir, _ := configDir()

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			fmt.Println("Unable to create config : ", err.Error())
			os.Exit(1)
		}
	}
}

func configDir() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = ".config"
	}
	defaultDir, err := defaultDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(defaultDir, configDir, DefaultDir), nil
}

func defaultDirPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

// parseYaml performs the real YAML parsing.
func parseYaml(confFile string) (*Config, error) {
	cfg, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	t := Config{}

	err = yaml.Unmarshal(cfg, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
