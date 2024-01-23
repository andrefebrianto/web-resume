package config

import "time"

type Application struct {
	Service Service `mapstructure:"service"`
}

type Service struct {
	PostgreSQL PostgreSQL `mapstructure:"postgresql"`
}

type PostgreSQL struct {
	ConnConfig PostgreSQLConnConfig `mapstructure:"connection-config"`

	Credential          DatabaseCredential `mapstructure:"credential"`
	MigrationCredential DatabaseCredential `mapstructure:"migration-credential"`

	Primary PostgreSQLInstance `mapstructure:"primary"`
}

type PostgreSQLConnConfig struct {
	MaxOpen     int32         `mapstructure:"maxopen"`
	MaxIdle     int32         `mapstructure:"maxidle"`
	MaxIdleTime time.Duration `mapstructure:"maxidletime"`
}

type PostgreSQLInstance struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	DBName string `mapstructure:"database"`
}

type DatabaseCredential struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
