package postgresql

const defaultPort = 5432

type Config struct {
	Database string `config:"DATABASE_NAME" yaml:"database"`
	User     string `config:"DATABASE_USER" yaml:"user"`
	Password string `config:"DATABASE_PASSWORD" yaml:"password"`
	Host     string `config:"DATABASE_HOST" yaml:"host"`
	Port     int    `config:"DATABASE_PORT" yaml:"port"`
	Retries  int    `config:"DB_CONNECT_RETRY" yaml:"retries"`
	PoolSize int    `config:"DB_POOL_SIZE" yaml:"pool_size"`
}

func (c *Config) withDefaults() (conf Config) {
	if c != nil {
		conf = *c
	}
	if conf.Port == 0 {
		conf.Port = defaultPort
	}
	return conf
}

func (c *Config) ReturnDatabase() string {
	return c.Database
}

func (c *Config) ReturnUser() string {
	return c.User
}

func (c *Config) ReturnPassword() string {
	return c.Password
}

func (c *Config) ReturnHost() string {
	return c.Host
}

func (c *Config) ReturnPort() int {
	if c.Port == 0 {
		return defaultPort
	}
	return c.Port
}

func (c *Config) ReturnRetries() int {
	return c.Retries
}

func (c *Config) ReturnPoolSize() int {
	return c.PoolSize
}
