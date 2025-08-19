package config

import (
	"log"
	"strings"
	"github.com/spf13/viper"
)

type ConfigStruct struct {
	Port   string `mapstructure:"PORT"`
	Env    string `mapstructure:"ENV"`

	
	DbUrl string `mapstructure:"DATABASE_URL"`

	
	JwtAccessSecret             string `mapstructure:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret            string `mapstructure:"JWT_REFRESH_SECRET"`
	AdminAccessTokenExpiresDays int    `mapstructure:"ADMIN_ACCESS_TOKEN_EXPIRES_DAYS"`
	AdminRefreshTokenExpiresDays int   `mapstructure:"ADMIN_REFRESH_TOKEN_EXPIRES_DAYS"`
	RefreshTokenExpiresDays     int    `mapstructure:"REFRESH_TOKEN_EXPIRES_DAYS"`
	AccessTokenExpiresDays      int    `mapstructure:"ACCESS_TOKEN_EXPIRES_DAYS"`

	// Lenhub
	LenhubClientID string `mapstructure:"LENHUB_CLIENT_ID"`
	LenhubApiKey   string `mapstructure:"LENHUB_API_KEY"`
	LenhubSenderID string `mapstructure:"LENHUB_SENDER_ID"`
	LenhubApiBase  string `mapstructure:"LENHUB_API_BASE"`

	// BuzzyCash
	BuzzyCashUsername  string `mapstructure:"BUZZY_CASH_USERNAME"`
	BuzzyCashPassword  string `mapstructure:"BUZZY_CASH_PASSWORD"`
	BuzzyCashCompanyID string `mapstructure:"BUZZY_CASH_COMPANYID"`
	BuzzyCashSenderID  string `mapstructure:"BUZZYCASH_SENDER_ID"`

	// Maekandex Gaming
	MaekandexGamingUrl string `mapstructure:"MAEKANDEX_GAMING_URL"`

	// Hubtel
	HubtelClientID     string `mapstructure:"HUBTEL_CLIENT_ID"`
	HubtelClientSecret string `mapstructure:"HUBTEL_CLIENT_SECRET"`
	HubtelSenderID     string `mapstructure:"HUBTEL_SENDER_ID"`
	HubtelApiBase      string `mapstructure:"HUBTEL_API_BASE"`

	// Super Admin
	SuperAdminPass string `mapstructure:"SUPER_ADMIN_PASS"`
}

var AppConfig ConfigStruct

func LoadConfig() {
	viper.SetConfigFile(".env")
	
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	



	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}
	if AppConfig.DbUrl == "" {
			log.Fatalf("❌ DATABASE_URL missing! (env not passed or mismatched key)")
		}
}
