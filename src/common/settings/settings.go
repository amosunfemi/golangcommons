package settings

import (
	"common/cfg"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var environments = map[string]string{
	"production":    cfg.APP_DIR + "security/settings/prod.json",
	"preproduction": cfg.APP_DIR + "security/settings/pre.json",
	"tests":         cfg.APP_DIR + "security/settings/tests.json",
}

type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

var settings Settings = Settings{}
var env = "preproduction"

func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(cfg.APP_DIR + environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

func GetEnvironment() string {
	return env
}

func Get() Settings {

	if &settings == nil || settings.PrivateKeyPath == "" {
		Init()
	}
	return settings
}

func IsTestEnvironment() bool {
	return env == "tests"
}
