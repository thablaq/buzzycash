package config

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ConfigStruct struct {
	Port string `envconfig:"PORT" default:"5005"`
	Env  string `envconfig:"env" default:"development"`
	
	DbUrl string `envconfig:"DATABASE_URL" required:"true"`
	
	JwtAccessSecret             string `envconfig:"JWT_ACCESS_SECRET" required:"true"`
	JwtRefreshSecret            string `envconfig:"JWT_REFRESH_SECRET" required:"true"`
	AdminAccessTokenExpiresDays int    `envconfig:"ADMIN_ACCESS_TOKEN_EXPIRES_DAYS"`
	AdminRefreshTokenExpiresDays int   `envconfig:"ADMIN_REFRESH_TOKEN_EXPIRES_DAYS"`
	RefreshTokenExpiresDays     int    `envconfig:"REFRESH_TOKEN_EXPIRES_DAYS"`
	AccessTokenExpiresDays      int    `envconfig:"ACCESS_TOKEN_EXPIRES_DAYS"`
	
	// Lenhub
	LenhubClientID string `envconfig:"LENHUB_CLIENT_ID"`
	LenhubApiKey   string `envconfig:"LENHUB_API_KEY"`
	LenhubSenderID string `envconfig:"LENHUB_SENDER_ID"`
	LenhubApiBase  string `envconfig:"LENHUB_API_BASE"`
	
	// BuzzyCash
	BuzzyCashUsername  string `envconfig:"BUZZY_CASH_USERNAME"`
	BuzzyCashPassword  string `envconfig:"BUZZY_CASH_PASSWORD"`
	BuzzyCashCompanyID string `envconfig:"BUZZY_CASH_COMPANYID"`
	BuzzyCashSenderID  string `envconfig:"BUZZYCASH_SENDER_ID"`
	
	// Maekandex Gaming
	MaekandexGamingUrl string `envconfig:"MAEKANDEX_GAMING_URL"`
	
	// Hubtel
	HubtelClientID     string `envconfig:"HUBTEL_CLIENT_ID"`
	HubtelClientSecret string `envconfig:"HUBTEL_CLIENT_SECRET"`
	HubtelSenderID     string `envconfig:"HUBTEL_SENDER_ID"`
	HubtelApiBase      string `envconfig:"HUBTEL_API_BASE"`
	
	//Flutterwave
	FlutterwaveSecretKey string `envconfig:"FLUTTERWAVE_SECRET_KEY"`
	FlutterwavePublicKey string `envconfig:"FLUTTERWAVE_PUBLIC_KEY"`
	FlutterwaveApiBase   string `envconfig:"FLUTTERWAVE_BASE_URL" default:"https://api.flutterwave.com/v3/"`
	FlutterwaveHashKey string `envconfig:"FLUTTERWAVE_HASH_KEY"`
	
	// Super Admin
	SuperAdminPass string `envconfig:"SUPER_ADMIN_PASS"`
}

var AppConfig ConfigStruct

func LoadConfig() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	

	if err := envconfig.Process("", &AppConfig); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}
	
	log.Println("✅ Configuration loaded successfully")
}


