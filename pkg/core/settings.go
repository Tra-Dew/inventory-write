package core

// Settings ...
type Settings struct {
	Port int32 `yaml:"port"`
	JWT  *JWT  `yaml:"jwt"`
}

// JWT ...
type JWT struct {
	Secret string `yaml:"secret"`
}
