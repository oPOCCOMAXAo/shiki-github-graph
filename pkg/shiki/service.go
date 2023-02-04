package shiki

type Service struct {
	config Config
	*api
}

type Config struct{}

func New(config Config) (*Service, error) {
	return &Service{
		config: config,
		api:    newAPI(),
	}, nil
}
