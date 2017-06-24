package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"mmsghana/cfg"
	"mmsghana/model"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	//"github.com/grafana/grafana/pkg/models"

	settings "mmsghana/settings"

	"mmsghana/redis"

	"golang.org/x/crypto/bcrypt"
)

//JWTAuthenticationBackend ...
type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var (
	setting_dir = ""
)

var authBackendInstance *JWTAuthenticationBackend = nil

// InitJWTAuthenticationBackend ...
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {

	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

// GenerateToken ...
func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string, userName string) (string, error) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	//expireCookie := time.Now().Add(time.Hour * 24)

	claims := MyCustomClaims{
		userName,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Subject:   userUUID,
			ExpiresAt: expireToken,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	/*token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
		"iat": time.Now().Unix(),
		"sub": userUUID,
	}*/

	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

//Authenticate ...
func (backend *JWTAuthenticationBackend) Authenticate(user model.User, password string) bool {
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	/*testUser := util.User{
		UUID:     uuid.New(),
		Username: user.Username,
		Password: string(hashedPassword),
	}*/
	//user.UUID = uuid.New()
	//fmt.Println(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(hashedPassword)))
	fmt.Println(password)
	fmt.Println(user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	claims := token.Claims.(jwt.MapClaims)
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(claims["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)
	fmt.Println(redisConn)
	if redisToken == nil {
		return false
	}

	return true
}

func getPrivateKey() *rsa.PrivateKey {
	fmt.Println(settings.Get().PrivateKeyPath)
	privateKeyFile, err := os.Open(cfg.APP_DIR + "security/settings/keys/private_key") ///Users/amosunsunday/Documents/Official/xtremepay/Sources/services/src/xtremepay.com/backoffice/security/settings/
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	//settings.Get().PublicKeyPath
	///Users/amosunsunday/Documents/Official/xtremepay/Sources/services/src/xtremepay.com/backoffice/security/settings/keys/public_key.pub
	publicKeyFile, err := os.Open(cfg.APP_DIR + "security/settings/keys/public_key.pub")
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
