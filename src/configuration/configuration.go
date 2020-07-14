package configuration

import (
	"fmt"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"
)

var Config *Configuration

//Configuration config
type Configuration struct {
	DBName string
	DBHost string
	DBPort int
	DBUser string
	DBPass string
}

//InitConfiguration init
func InitConfiguration() {
	port, _ := strconv.Atoi(GetEnv("DB_PORT", "3306"))
	Config = &Configuration{
		DBName: GetEnv("DB_NAME", "shop"),
		DBUser: GetEnv("DB_USER", "root"),
		DBPass: GetEnv("DB_PASS", ""),
		DBPort: port,
		DBHost: GetEnv("DB_HOST", "localhost"),
	}
}

//GetConnectionString get connection string
func (c *Configuration) GetConnectionString() string {
	log.Error("test")
	if c.DBHost == "" || c.DBPass == "" {
		InitConfiguration()
	}

	if c.DBPass == "" {
		return fmt.Sprintf("%s@%s:%d/%s?charset=utf8&parseTime=True&loc=Local", c.DBUser, c.DBHost, c.DBPort, c.DBName)
	}

	return fmt.Sprintf("%s:%s@%s:%d/%s?charset=utf8&parseTime=True&loc=Local", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

//GetEnv get env
func GetEnv(v string, d string) string {
	val := os.Getenv(v)

	if val == "" {
		return d
	}

	return val
}
