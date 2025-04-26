package conf

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/webbythien/monorepo/sdk/api/server"
)

type NextBackend struct {
	Address string `yaml:"address" mapstructure:"address"`
}
type Temporal struct {
	Address   string `yaml:"address" mapstructure:"address"`
	Namespace string `yaml:"namespace" mapstructure:"namespace"`
}
type Auth struct {
	Secret string `yaml:"secret" mapstructure:"secret"`
}
type General struct {
	Company string `yaml:"company" mapstructure:"company"`
	Domain  string `yaml:"domain" mapstructure:"domain"`
}

type S3 struct {
	Bucket string `yaml:"bucket" mapstructure:"bucket"`
	Region string `yaml:"region" mapstructure:"region"`
	SqsURL string `yaml:"sqs_url" mapstructure:"sqs_url"`
}

// PostgreSQL is settings of a PostgreSQL server. It contains almost same fields as postgresql.Config,
// but with some different field names and tags.
type PostgreSQL struct {
	Host        string `yaml:"host" mapstructure:"host"`
	Port        string `yaml:"port" mapstructure:"port"`
	Database    string `yaml:"database" mapstructure:"database"`
	Username    string `yaml:"username" mapstructure:"username"`
	Password    string `yaml:"password" mapstructure:"password"`
	SSLMode     string `yaml:"sslmode" mapstructure:"sslmode"`
	SSLRootCert string `yaml:"sslrootcert" mapstructure:"sslrootcert"`
	SearchPath  string `yaml:"search_path" mapstructure:"search_path"`
	RDS         *RDS   `yaml:"rds" mapstructure:"rds"`
}

// FormatDSN returns PostgreSQL DSN from settings.
func (m *PostgreSQL) FormatDSN() string {
	_, err := strconv.Atoi(m.Port)
	if err != nil {
		panic(fmt.Errorf("invalid port: %w", err))
	}
	var dsn string
	for key, val := range map[string]string{
		"host":        m.Host,
		"port":        m.Port,
		"user":        m.Username,
		"password":    m.Password,
		"dbname":      m.Database,
		"sslmode":     m.SSLMode,
		"sslrootcert": m.SSLRootCert,
		"search_path": m.SearchPath,
	} {
		if val != "" {
			dsn += fmt.Sprintf(" %s=%s", key, val)
		}
	}
	dsn = strings.TrimLeft(dsn, " ")
	return dsn
}

type RDS struct {
	Region    string `yaml:"region" mapstructure:"region"`
	IsEnabled bool   `yaml:"is_enabled" mapstructure:"is_enabled"`
}

// MySQL is settings of a MySQL server. It contains almost same fields as postgresql.Config,
// but with some different field names and tags.
type MySQL struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Database string `yaml:"database" mapstructure:"database"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
}

// FormatDSN returns MySQL DSN from settings.
func (m *MySQL) FormatDSN() string {
	port, err := strconv.Atoi(m.Port)
	if err != nil {
		panic(fmt.Errorf("invalid port: %w", err))
	}
	cfg := mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", m.Host, port),
		DBName:               m.Database,
		User:                 m.Username,
		Passwd:               m.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	return cfg.FormatDSN()
}

type Mongo struct {
	URI      string `yaml:"uri" mapstructure:"uri"`
	Database string `yaml:"database" mapstructure:"database"`
}

type Base struct {
	Server      server.Config `yaml:"server" mapstructure:"server"`
	ServiceName string        `yaml:"service_name" mapstructure:"service_name"`
	Environment string        `yaml:"environment" mapstructure:"environment"`
}

type Kafka struct {
	Host         string   `yaml:"host" mapstructure:"host"`
	Port         int      `yaml:"port" mapstructure:"port"`
	Brokers      string   `yaml:"brokers" mapstructure:"brokers"`
	GroupID      string   `yaml:"group_id" mapstructure:"group_id"`
	GroupTopics  []string `yaml:"group_topics" mapstructure:"group_topics"`
	Topic        string   `yaml:"topic" mapstructure:"topic"`
	SASLProtocol string   `yaml:"sasl_protocol" mapstructure:"sasl_protocol"`
	User         string   `yaml:"user" mapstructure:"user"`
	Password     string   `yaml:"password" mapstructure:"password"`
	Cert         string   `yaml:"cert" mapstructure:"cert"`
	Key          string   `yaml:"key" mapstructure:"key"`
}

func (k *Kafka) BrokersList() []string {
	return strings.Split(k.Brokers, ",")
}

func (k *Kafka) MustLoadCert() {
	if k.SASLProtocol != "scram" {
		return
	}

	if k.Cert == "" {
		cert, err := os.ReadFile("server.crt")
		if err != nil {
			panic(err)
		}
		k.Cert = string(cert)
	}

	if k.Key == "" {
		key, err := os.ReadFile("server.key")
		if err != nil {
			panic(err)
		}
		k.Key = string(key)
	}
}

type Firebase struct {
	ProjectID string `yaml:"project_id" mapstructure:"project_id"`
}

type Cloudflare struct {
	APIToken  string `yaml:"api_token" mapstructure:"api_token"`
	AccountID string `yaml:"account_id" mapstructure:"account_id"`
}

type Observe struct {
	Metric struct {
		Enabled          bool   `yaml:"enabled" mapstructure:"enabled"`
		EndpointExporter string `yaml:"endpoint_exporter" mapstructure:"endpoint_exporter"`
	} `yaml:"metric" mapstructure:"metric"`
	Trace struct {
		Enabled          bool    `yaml:"enabled" mapstructure:"enabled"`
		EndpointExporter string  `yaml:"endpoint_exporter" mapstructure:"endpoint_exporter"`
		SampleRate       float64 `yaml:"sample_rate" mapstructure:"sample_rate"`
	} `yaml:"trace" mapstructure:"trace"`
}
