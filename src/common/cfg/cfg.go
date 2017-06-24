package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Unique identifier for cookie in this app
var COOKIE_NAME string = "GONGFU-COOKIE"
var HASHKEY string = "gongfu-hash-key"
var BLOCKKEY string = "gongfu-block-key"

// defaults
var APPNAME string = "gongfu"
var DB_SERVER string = "localhost"
var DB_TYPE string = "postgres"
var DB_NAME string = "gongfu"
var DB_USER string = "gongfu"
var DB_PASSWORD string = "gongfu"
var HTTP_PORT string = ""
var SERVER string = "gongfu"
var APP_DIR string = ""
var RPT_SERVER_URL string = ""
var ADCLUSTERPATH string = ""
var SOVPATH string = ""
var DB_PORT int64 = 5432

var conf JsonCfg

type JsonCfg struct {
	Appname         string
	Http_port       string
	Server          string
	Db_server       string
	Db_name         string
	Db_user         string
	Db_password     string
	Application_dir string
	Rptserverurl    string
	Adcluster_path  string
	Sov_path        string
}

// ConfigFrom is a method to read json properties
func (cfg *JsonCfg) ConfigFrom(path string) (err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &cfg)
	if err != nil {
		log.Printf("bad json", err)
	}

	return
}

// ConfigSetup is the init function called from app.go, to read the cfg.json file
// If error reading the cfg.json; defaults are used
func Setup(config map[string]interface{}) {
	fmt.Println(config)
	APPNAME = config["appname"].(string) //conf.Appname
	APP_DIR = config["appdir"].(string)
	HTTP_PORT = fmt.Sprintf("%s", config["http_port"].(string))

}
