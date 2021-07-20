package core

// Settings ...
type Settings struct {
	Port   int32          `yaml:"port"`
	JWT    *JWT           `yaml:"jwt"`
	SQS    *SessionConfig `yaml:"sqs"`
	SNS    *SessionConfig `yaml:"sns"`
	Events *Events        `yaml:"events"`
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

// Events ...
type Events struct {
	ItemsLockRequested string `yaml:"items-lock-requested"`
	ItemsLockCompleted string `yaml:"items-lock-completed"`
}
