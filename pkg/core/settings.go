package core

// Settings ...
type Settings struct {
	Port     int32           `yaml:"port"`
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
	ItemsCreated       string `yaml:"items-created"`
	ItemsUpdated       string `yaml:"items-updated"`
	ItemsLockRequested string `yaml:"items-lock-requested"`
	ItemsLockCompleted string `yaml:"items-lock-completed"`
}
