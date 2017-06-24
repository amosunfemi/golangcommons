package utility

import (
	"common/cfg"
	"common/sys"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	"github.com/knq/dburl"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//DB ...
var DB gorm.DB
var DBSql sql.DB

//DataStore ...
var DataStoreObject *DataStore

//Keebler ...
var Keebler *securecookie.SecureCookie

//SetupDB ...
func SetupDB(config map[string]interface{}, entites []interface{}, dbtype string) {
	cfg.DB_SERVER = config["address"].(string)
	cfg.DB_USER = config["userid"].(string)
	cfg.DB_TYPE = config["dbtype"].(string)
	cfg.DB_NAME = config["dbname"].(string)
	cfg.DB_PASSWORD = config["password"].(string)
	cfg.DB_PORT = int64(config["port"].(float64))
	/*settings := "user=" + cfg.DB_USER + " password=" + cfg.DB_PASSWORD + " dbname=" + cfg.DB_NAME + " sslmode=disable"
	db, err := sql.Open("postgres", settings)

	PanicIf(err)
	DB = db
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	log.Printf("DB has been pinged")*/
	connString := fmt.Sprintf("host=%s user=%s port=%d password=%s dbname=%s sslmode=disable", cfg.DB_SERVER, cfg.DB_USER, cfg.DB_PORT, cfg.DB_PASSWORD, cfg.DB_NAME)
	db, err := gorm.Open("postgres", connString)
	DataStoreObject = &DataStore{}
	DataStoreObject.StoreType = dbtype
	DB.LogMode(true)
	DB = *db
	DataStoreObject.RDBMS = RDBMSImpl{db}

	DataStoreObject.InitSchema(entites)
	PanicIf(err)

	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	log.Printf("DB has been pinged")

}

//SetupDB ...
func SetupDB2(config map[string]interface{}) {
	cfg.DB_SERVER = config["address"].(string)
	cfg.DB_USER = config["userid"].(string)
	cfg.DB_TYPE = config["dbtype"].(string)
	cfg.DB_NAME = config["dbname"].(string)
	cfg.DB_PASSWORD = config["password"].(string)
	cfg.DB_PORT = int64(config["port"].(float64))
	/*settings := "user=" + cfg.DB_USER + " password=" + cfg.DB_PASSWORD + " dbname=" + cfg.DB_NAME + " sslmode=disable"
	db, err := sql.Open("postgres", settings)

	PanicIf(err)
	DB = db
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	log.Printf("DB has been pinged")*/
	connString := fmt.Sprintf("pgsql://u%s:%s@%s:%d/%s?sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_SERVER, cfg.DB_PORT, cfg.DB_NAME)
	db, err := dburl.Open(connString)
	PanicIf(err)
	DBSql = *db
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	log.Printf("DB has been pinged")
	//DBSql.LogMode(true)
}

//SetupKeebler ...
func SetupKeebler() {
	Keebler = securecookie.New([]byte(cfg.HASHKEY), []byte(cfg.BLOCKKEY))
	log.Printf("Keebler has been created.... ready to make cookies...")
}

//GenHash ...
func GenHash(c int) string {
	//c := 8
	b := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

//CookieHasAccess ...
func CookieHasAccess(c *gin.Context) (bool, string, string) {
	if cookie, err := c.Request.Cookie(cfg.COOKIE_NAME); err == nil {
		value := make(map[string]string)
		if err = Keebler.Decode(cfg.COOKIE_NAME, cookie.Value, &value); err == nil {
			log.Printf("The value of access is " + value[sys.COOKIE_APP_ACCESS])
			log.Printf("The value of email is " + value[sys.COOKIE_EMAIL])
			log.Printf("Access value of [COOKIE_APP_ACCESS]: %v", value[sys.COOKIE_APP_ACCESS])
			if value[sys.COOKIE_APP_ACCESS] == sys.ACCESS_OK {
				return true, value[sys.COOKIE_EMAIL], value[sys.COOKIE_USERID]
			}
		}
	}
	log.Printf("Where is the cookie? %v", cfg.COOKIE_NAME)
	return false, "", ""
}

//PanicIf ...
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// Schedule is a generic timer func that will run a func after a delay
func Schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
	// stop is the channel that will stop if you do: stop <- true
}

//HashPassword ...
func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash, err
}

//Round ...
func Round(val float64, prec int) float64 {
	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(prec))
}

// Btc_Sat converts BTC float to Satochi Int; 100 million sat = 1BTC
// int64 value +/- 9223372036854775807
//Btc_Sat ...
func Btc_Sat(btc float64) (sat int64, err error) {
	num := btc * 100000000
	if (num > 9223372036854775807) || (num < -9223372036854775807) {
		num = -1
		err = errors.New("Out of scale for int64")
	}
	return int64(num), err
}

// TODO write check for float64 size
func Sat_Btc(sat int64) (btc float64, err error) {
	num := float64(sat)
	if (num > 9223372036854775807) || (num < -9223372036854775807) {
		num = -1
		err = errors.New("Out of scale for float64")
	}
	btc = num / 100000000
	return btc, err
}

//RandStr ...
func RandStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
