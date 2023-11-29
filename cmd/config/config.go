package config

const (
	Debug = "Debug"
	Prod  = "Prod"
)

type Config struct {
	Server   Server   `yaml:"Server"`
	Postgers Postgers `yaml:"Postgres"`
	Redis    Redis    `yaml:"Redis"`
}

type (
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Postgers struct {
		DbName   string `yaml:"db_name"`
		DbUser   string `yaml:"db_user"`
		Port     string `yaml:"db_port"`
		Password string `yaml:"db_password"`
		SslMode  bool   `yaml:"sslmode"`
	}

	Redis struct {
		Addres string `yaml:"addres"`
	}
)
