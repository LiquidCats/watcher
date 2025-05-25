package docker

import (
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/go-faster/errors"
)

const prefix = "/run/secrets/"

func GetSecret(name string) (string, error) {
	if strings.HasPrefix(name, prefix) {
		file, err := os.OpenFile(name, os.O_RDONLY, 0600)
		if err != nil {
			return "", errors.Wrap(err, "failed to open secret file")
		}
		defer file.Close()

		var buff []byte

		_, err = io.ReadFull(file, buff)
		if err != nil {
			return "", errors.Wrap(err, "failed to read secret file")
		}

		var decodedBuff []byte

		_, err = base64.StdEncoding.Decode(decodedBuff, buff)
		if err != nil {
			return "", errors.Wrap(err, "failed to decode secret")
		}

		return string(decodedBuff), nil
	}

	return name, nil
}
