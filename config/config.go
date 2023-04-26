package config

import (
	"errors"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var (
	cfg *Config
)

type Config struct {
	Server         ServerConfig
	Postgres       PostgresConfig
	Redis          RedisConfig
	Jwt            JwtConfig
	FirstSuperUser FirstSuperUserConfig
	Logger         Logger
	SmtpEmail      SmtpEmailConfig
	Email          EmailConfig
	TaskRedis      TaskRedisConfig
}

type ServerConfig struct {
	AppVersion     string
	Port           string
	Mode           string
	ProcessTimeout int
	ReadTimeout    int
	WriteTimeout   int
	MigrateOnStart bool
}

type Logger struct {
	Encoding string
	Level    string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type RedisConfig struct {
	Addr         string
	Password     string
	Db           int
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
}

type TaskRedisConfig struct {
	Addr string
	Db   int
}

type JwtConfig struct {
	SecretKey                  string
	Issuer                     string
	AccessTokenExpireDuration  int64
	AccessTokenPrivateKey      string
	AccessTokenPublicKey       string
	RefreshTokenExpireDuration int64
	RefreshTokenPrivateKey     string
	RefreshTokenPublicKey      string
}

type FirstSuperUserConfig struct {
	Email    string
	Name     string
	Password string
}

type EmailConfig struct {
	From                string
	Name                string
	Link                string
	LogoLink            string
	Copyright           string
	VerificationSubject string
	ResetSubject        string
}

type SmtpEmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	UseTls   bool
	UseSsl   bool
}

func ToSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func BindEnvs(vp *viper.Viper, iface interface{}, partsKey []string, partsEnvKey []string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)

		tv := strings.ToUpper(ToSnakeCase(t.Name))

		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(vp, v.Interface(), append(partsKey, t.Name), append(partsKey, tv))
		default:
			key := strings.ToLower(strings.Join(append(partsKey, t.Name), "."))
			envKey := strings.ToUpper(strings.Join(append(partsEnvKey, tv), "_"))

			vp.BindEnv(key, envKey) //nolint:errcheck
		}
	}
}

// Load config file from given path
func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("config/config.default")
	v.SetConfigType("yml")

	BindEnvs(v, Config{}, []string{}, []string{})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	cfg = &c

	return &c, nil
}

func GetCfg() *Config {
	if cfg == nil {
		cfg = new(Config)
	}
	return cfg
}
