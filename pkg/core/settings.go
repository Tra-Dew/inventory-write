package core

// Settings ...
type Settings struct {
	Port     int32           `yaml:"port"`
	GRPCPort int32           `yaml:"grpc_port"`
	JWT      *JWT            `yaml:"jwt"`
	SQS      *SessionConfig  `yaml:"sqs"`
	SNS      *SessionConfig  `yaml:"sns"`
	Postgres *PostgresConfig `yaml:"postgres"`
	Events   *Events         `yaml:"events"`
}

// JWT ...
type JWT struct {
	Secret string `yaml:"secret"`
}

// SessionConfig ...
type SessionConfig struct {
	Region   string `yaml:"region"`
	Endpoint string `yaml:"endpoint"`
	Path     string `yaml:"path"`
	Profile  string `yaml:"profile"`
	Fake     bool   `yaml:"fake"`
}

// PostgresConfig ...
type PostgresConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
}

// Events ...
type Events struct {
	ItemsUpdated string `yaml:"items-updated"`
}
