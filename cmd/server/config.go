package main

import "market/pkg/storage/postgresql"

type Config struct {
	Host       string `config:"HOST" yaml:"host"`
	Port       string `config:"PORT" yaml:"port"`
	Migrations bool   `config:"MIGRATIONS" yaml:"migrations"`

	DatabaseType   string            `config:"DATABASE_TYPE" yaml:"database_type"`
	DatabaseConfig postgresql.Config `config:"DATABASE_CONFIG" yaml:"database_config"`
}
