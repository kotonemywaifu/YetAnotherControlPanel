package others

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

type Config struct {
	Port            int       `json:"port"`
	Password        string    `json:"password"`
	TrustedHosts    []string  `json:"trusted-hosts"`
	SecuredEntrance string    `json:"secured-entrance"` // this can make users get away from password brute force attack
	Session         string    `json:"session"`          // login token in cookie
	Gin             ConfigGin `json:"gin"`
}

type ConfigGin struct {
	DebugMode bool `json:"debug-mode"`
}

var configDir string
var TheConfig *Config

var SESSION_TOKEN_LENGTH = 128

func makeConfig() *Config {
	conf := new(Config)

	// initialize values
	conf.Port = 8080
	conf.Password = util.RandStringRunes(10)
	conf.TrustedHosts = []string{}
	conf.SecuredEntrance = "" /* util.randomString(8) */
	conf.Session = util.RandStringRunes(SESSION_TOKEN_LENGTH)

	conf.Gin.DebugMode = false

	return conf
}

func InitEnv() error {
	// random in golang seed is not random, so we need to use time.Now()
	rand.Seed(time.Now().UnixNano())

	return nil
}

func InitConfig() error {
	TheConfig = makeConfig()

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
