package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHttpPort               = "8080"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
	defaultDataFilePathSms        = "/var/data/sms.data"
	defaultDataFilePathEmail      = "/var/data/email.data"
	defaultDataFilePathVoiceCall  = "/var/data/voice.data"
	defaultDataFilePathBilling    = "/var/data/billing.data"
	defaultDataUrlMms             = "http://localhost:8383/mms"
	defaultDataUrlIncident        = "http://localhost:8383/accendent"
	defaultDataUrlSupport         = "http://localhost:8383/support"
	defaultCacheExpiration        = 30 * time.Second
	defaultCacheCleanupInterval   = 10 * time.Second
)

type (
	Config struct {
		HTTP         HTTPConfig
		DataUrl      DataUrlConfig
		DataFilePath DataFilePathConfig
		Cache        CacheConfig
	}

	DataFilePathConfig struct {
		Sms       string `mapstructure:"sms"`
		VoiceCall string `mapstructure:"voicecall"`
		Email     string `mapstructure:"email"`
		Billing   string `mapstructure:"billing"`
	}

	DataUrlConfig struct {
		Mms      string `mapstructure:"mms"`
		Incident string `mapstructure:"incident"`
		Support  string `mapstructure:"support"`
	}

	HTTPConfig struct {
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	CacheConfig struct {
		Expiration      time.Duration `mapstructure:"expiration"`
		CleanupInterval time.Duration `mapstructure:"cleanupInterval"`
	}
)

func Init(path string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	if err := parseEnv(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)
	fmt.Println(cfg)
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("dataFilePath", &cfg.DataFilePath); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("dataUrl", &cfg.DataUrl); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("cache", &cfg.Cache); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0]) // folder
	viper.SetConfigName(path[1]) // config file name

	return viper.ReadInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.max_header_megabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("dataFilePath.sms", defaultDataFilePathSms)
	viper.SetDefault("dataFilePath.voicecall", defaultDataFilePathVoiceCall)
	viper.SetDefault("dataFilePath.email", defaultDataFilePathEmail)
	viper.SetDefault("dataFilePath.billing", defaultDataFilePathBilling)
	viper.SetDefault("dataUrl.mms", defaultDataUrlMms)
	viper.SetDefault("dataUrl.incident", defaultDataUrlIncident)
	viper.SetDefault("dataUrl.support", defaultDataUrlSupport)
	viper.SetDefault("cache.expiration", defaultCacheExpiration)
	viper.SetDefault("cache.cleanupInterval", defaultCacheCleanupInterval)
}

func setFromEnv(cfg *Config) {
	cfg.HTTP.Port = viper.GetString("PORT")
}

func parseEnv() error {
	return parseHerokuEnvVariables()
}

func parseHerokuEnvVariables() error {
	if err := viper.BindEnv("PORT"); err != nil {
		return err
	}
	return nil
}
