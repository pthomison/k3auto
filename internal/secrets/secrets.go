package secrets

import "context"

type SecretConfig struct {
	DefaultName      string   `yaml:"DefaultSecret,omitempty"`
	DefaultNamespace string   `yaml:"DefaultNamespace,omitempty"`
	Secrets          []Secret `yaml:"Secrets,omitempty"`
}

func NewSecretConfig(defaultCapacity int, defaultSecretCapacity int) SecretConfig {
	return SecretConfig{
		Secrets: func() []Secret {
			s := []Secret{}
			for range defaultCapacity {
				s = append(s, NewSecret(defaultSecretCapacity))
			}
			return s
		}(),
	}
}

type Secret struct {
	Type            string   `yaml:"Type,omitempty"`
	Args            []string `yaml:"Args,omitempty"`
	SecretName      string   `yaml:"SecretName,omitempty"`
	SecretKey       string   `yaml:"SecretKey,omitempty"`
	SecretNamespace string   `yaml:"SecretNamespace,omitempty"`
}

func NewSecret(defaultCapacity int) Secret {
	return Secret{
		Args: make([]string, defaultCapacity),
	}
}

type SecretResolver interface {
	Resolve(ctx context.Context, args []string) (string, error)
}

var (
	ResolverMap = map[string]SecretResolver{
		"exec": &ExecResolver{},
		// "paramstore": &ParamStoreResolver{},
	}
)
