package configs

import (
	"fmt"
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module is the fx module that provides the configuration.
var Module = fx.Options(
	fx.Provide(ProvideConfig),
)

// Config holds the application wide configurations.
// The values are read by viper from the config file or environment variables.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	GRPC   GRPCConfig   `mapstructure:"grpc"`
	PSQL   PSQLConfig   `mapstructure:"postgres"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type GRPCConfig struct {
	UserHost string `mapstructure:"user_host"`
	UserPort int    `mapstructure:"user_port"`
	AuthHost string `mapstructure:"auth_host"`
	AuthPort int    `mapstructure:"auth_port"`
}

// PSQLConfig holds PostgreSQL connection configuration.
type PSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type JWTConfig struct {
	SecretKey            string        `mapstructure:"secret_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

// NewConfig creates a new Config instance.
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()

	// Set default values for the configuration.
	setDefaults(v)

	// Read from environment variables
	v.AutomaticEnv()

	// Try to read from config file, but continue if not found
	v.AddConfigPath(path)
	v.SetConfigName("config.example")
	v.SetConfigType("yaml")

	// Try to read config file, but log the error instead of failing
	if err := v.ReadInConfig(); err != nil {
		logger.Sugar().Warnf("could not read config file: %v", err)
	}

	// Unmarshal the configuration into the config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	// Validate essential configuration values
	if err := validateServerConfig(config.Server); err != nil {
		return nil, err
	}

	if err := validateGRPCConfig(config.GRPC); err != nil {
		return nil, err
	}

	if err := validatePSQLConfig(config.PSQL); err != nil {
		return nil, err
	}

	if err := validateJWTConfig(config.JWT); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets default configuration values in viper.
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", "3001")

	v.SetDefault("grpc.user_host", "localhost")
	v.SetDefault("grpc.user_port", "8080")
	v.SetDefault("grpc.auth_host", "localhost")
	v.SetDefault("grpc.auth_port", "8081")

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5433")
	v.SetDefault("postgres.user", "root")
	v.SetDefault("postgres.password", "secret")
	v.SetDefault("postgres.database", "auth-db")
	v.SetDefault("postgres.ssl_mode", "disable")

	v.SetDefault("jwt.secret_key", "abcd1234abcd1234abcd1234")
	v.SetDefault("jwt.access_token_duration", "15m")
	v.SetDefault("jwt.refresh_token_duration", "24h")
}

// validateServerConfig ensures that essential server config values are present.
func validateServerConfig(serverConfig ServerConfig) error {
	if serverConfig.Port == 0 {
		return fmt.Errorf("server port is required")
	}
	return nil
}

// validateGRPCConfig ensures that essential gRPC config values are present.
func validateGRPCConfig(grpcConfig GRPCConfig) error {
	if grpcConfig.UserHost == "" {
		return fmt.Errorf("user grpc host is required")
	}
	if grpcConfig.UserPort == 0 {
		return fmt.Errorf("user grpc port is required")
	}
	if grpcConfig.AuthHost == "" {
		return fmt.Errorf("auth grpc host is required")
	}
	if grpcConfig.AuthPort == 0 {
		return fmt.Errorf("auth grpc port is required")
	}
	return nil
}

// validatePSQLConfig ensures that essential PSQL config values are present.
func validatePSQLConfig(psqlConfig PSQLConfig) error {
	if psqlConfig.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if psqlConfig.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if psqlConfig.User == "" {
		return fmt.Errorf("database user is required")
	}
	if psqlConfig.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if psqlConfig.Database == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

// validateJWTConfig ensures that essential JWT config values are present.
func validateJWTConfig(jwtConfig JWTConfig) error {
	if jwtConfig.SecretKey == "" {
		return fmt.Errorf("jwt secret key is required")
	}
	if jwtConfig.AccessTokenDuration == 0 {
		return fmt.Errorf("jwt access token duration is required")
	}
	if jwtConfig.RefreshTokenDuration == 0 {
		return fmt.Errorf("jwt refresh token duration is required")
	}
	return nil
}

// ProvideConfig is an fx provider that loads the configuration.
func ProvideConfig(logger *logger.Logger) (*Config, error) {
	return LoadConfig(".", logger)
}
