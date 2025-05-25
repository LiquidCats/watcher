package docker

import (
	"os"
	"strings"

	"github.com/go-faster/errors"
)

const prefix = "/run/secrets/"

func GetSecret(name string) (string, error) {
	if strings.HasPrefix(name, prefix) {
		data, err := os.ReadFile(name)
		if err != nil {
			return "", errors.Wrap(err, "failed to open secret file")
		}

		return strings.TrimSpace(string(data)), nil
	}

	return name, nil
}
