package secrets

import (
	"io"

	"gopkg.in/yaml.v3"
)

func LoadConfigFile(r io.Reader) (SecretConfig, error) {
	c := NewSecretConfig(10, 10)

	b, err := io.ReadAll(r)
	if err != nil {
		return SecretConfig{}, err
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return SecretConfig{}, err
	}

	return c, nil
}
