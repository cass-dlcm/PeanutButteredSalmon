package types

type Server struct {
	ShortName string `yaml:"short_name"`
	ApiKey string `yaml:"api_key"`
	Address string `yaml:"address"`
}