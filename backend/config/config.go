package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config is the top-level application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	Database DatabaseConfig `yaml:"database"`
	CORS     CORSConfig     `yaml:"cors"`
}

type ServerConfig struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Addr returns the listen address.
func (s ServerConfig) Addr() string {
	if p := os.Getenv("PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			s.Port = port
		}
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type JWTConfig struct {
	Secret  string `yaml:"secret"`
	Timeout string `yaml:"timeout"`
	Issuer  string `yaml:"issuer"`
}

type DatabaseConfig struct {
	Driver   string           `yaml:"driver"`
	SQLite   SQLiteConfig     `yaml:"sqlite"`
	MySQL    MySQLConfig      `yaml:"mysql"`
	Postgres PostgresConfig   `yaml:"postgres"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	ParseTime bool `yaml:"parseTime"`
}

func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database, m.Charset, m.ParseTime)
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.Username, p.Password, p.Database, p.SSLMode)
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allowOrigins"`
	AllowMethods []string `yaml:"allowMethods"`
	AllowHeaders []string `yaml:"allowHeaders"`
}

// global holds the loaded configuration.
var global *Config

// Load reads settings.yml (path overridable via CONFIG_PATH) into the global config.
func Load(path string) (*Config, error) {
	if path == "" {
		path = "config/settings.yml"
	}
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		path = env
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	global = &c
	return global, nil
}

// Get returns the global configuration (panics if not loaded).
func Get() *Config {
	if global == nil {
		panic("config not loaded; call config.Load first")
	}
	return global
}
