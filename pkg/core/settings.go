package core

// Settings ...
type Settings struct {
	Port int32 `yaml:"port"`
	JWT  *JWT  `yaml:"jwt"`
}

// JWT ...
type JWT struct {
	ExpirationInMinutes int64  `yaml:"expiration_in_minutes"`
	Secret              string `yaml:"jwt_secret"`
}
