package emitter

// EmitterConfig chứa các cấu hình cần thiết để kết nối với Emitter
type EmitterConfig struct {
	RedisAddr     string `yaml:"redis_addr" mapstructure:"redis_addr"`
	RedisPassword string `yaml:"redis_password" mapstructure:"redis_password"`
	RedisDB       int    `yaml:"redis_db" mapstructure:"redis_db"`
	TLSEnabled    bool   `yaml:"tls_enabled" mapstructure:"tls_enabled"`
	Prefix        string `yaml:"prefix" mapstructure:"prefix"`
	Uid           string `yaml:"uid" mapstructure:"uid"`
	Nsp           string `yaml:"nsp" mapstructure:"nsp"`
	EventType     int    `yaml:"event_type" mapstructure:"event_type"`
}
