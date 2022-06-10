package others

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port            int      `json:"port"`
	TrustedHosts    []string `json:"trusted-hosts"`
	Log             bool     `json:"log"`
	SecuredEntrance string   `json:"secured-entrance"` // this can make users get away from password brute force attack
	CacheTemplate   bool     `json:"cache-template"`   // lower CPU usage, but higher memory usage
	MinifyResources bool     `json:"minify-resources"` // minify resources, to reduce traffic
	FailToBan       struct {
		Enabled  bool `json:"enabled"`
		Failures int  `json:"failures"`
		BanTime  int  `json:"ban-time"` // in seconds
	} `json:"fail-to-ban"` // ban user if failed to login for x times
	Gin struct {
		DebugMode  bool `json:"debug-mode"`
		UseDebugFS bool `json:"use-debug-fs"`
		Tls        struct {
			Enabled  bool   `json:"enabled"`
			CertFile string `json:"cert"`
			KeyFile  string `json:"key"`
		} `json:"tls"`
	} `json:"gin"`
}

var ConfigDir string
var TheConfig *Config

func makeConfig() *Config {
	conf := new(Config)

	// initialize values
	conf.Port = 8080
	conf.TrustedHosts = []string{}
	conf.Log = true
	conf.SecuredEntrance = "" /* util.randomString(8) */
	conf.CacheTemplate = true
	conf.MinifyResources = false

	conf.FailToBan.Enabled = true
	conf.FailToBan.Failures = 5
	conf.FailToBan.BanTime = 180

	conf.Gin.DebugMode = false
	conf.Gin.UseDebugFS = false
	conf.Gin.Tls.Enabled = false
	conf.Gin.Tls.CertFile = ""
	conf.Gin.Tls.KeyFile = ""

	return conf
}

// func firstRun() error {

// }

func InitEnv() error {
	// random in golang seed is not random, so we need to use time.Now()
	rand.Seed(time.Now().UnixNano())

	// create directory if not exists
	osConfig, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	ConfigDir = osConfig + "/yacp/"

	// create config dir if not exists
	if _, err = os.Stat(ConfigDir); os.IsNotExist(err) {
		err = os.MkdirAll(ConfigDir, os.ModePerm)
		if err != nil {
			return err
		}
		// if yacp dir not exists, this is the first time to run yacp
		// we need to initialize something
		// err = firstRun()
		// if err != nil {
		// 	return err
		// }
	}

	subdirs := []string{"logs"}
	for _, subdir := range subdirs {
		if _, err = os.Stat(ConfigDir + subdir); os.IsNotExist(err) {
			err = os.MkdirAll(ConfigDir+subdir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InitConfig() error {
	TheConfig = makeConfig()

	return readConfigJson()
}

func readConfigJson() error {
	configFile := ConfigDir + "config.json"

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

	TheConfig.Gin.Tls.CertFile = strings.ReplaceAll(TheConfig.Gin.Tls.CertFile, "!config-dir!", ConfigDir[:len(ConfigDir)-1])
	TheConfig.Gin.Tls.KeyFile = strings.ReplaceAll(TheConfig.Gin.Tls.KeyFile, "!config-dir!", ConfigDir[:len(ConfigDir)-1])

	return writeConfigJson() // re-save config file to ensure data integrity
}

func writeConfigJson() error {
	configFile := ConfigDir + "config.json"

	file, err := json.MarshalIndent(TheConfig, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
