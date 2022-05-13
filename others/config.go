package others

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Port         int       `json:"port"`
	TrustedHosts []string  `json:"trusted-hosts"`
	Gin          ConfigGin `json:"gin"`
}

type ConfigGin struct {
	DebugMode bool `json:"debug-mode"`
}

var configDir string
var TheConfig = makeConfig()

func makeConfig() *Config {
	conf := new(Config)

	// initialize values
	conf.Port = 8080
	conf.TrustedHosts = []string{}

	conf.Gin.DebugMode = false

	return conf
}

func InitConfig() error {
	osConfig, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configDir = osConfig + "/yacp/"

	// create config dir if not exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return readConfigJson()
}

func readConfigJson() error {
	configFile := configDir + "config.json"

	// create config file if not exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("Config file not exists, creating...")
		err := writeConfigJson()
		if err != nil {
			return err
		}
		return nil
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, TheConfig)
	if err != nil {
		return err
	}

	return writeConfigJson() // re-save config file to ensure data integrity
}

func writeConfigJson() error {
	configFile := configDir + "config.json"

	file, err := json.MarshalIndent(TheConfig, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
