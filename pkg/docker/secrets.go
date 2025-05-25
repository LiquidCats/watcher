package docker

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-faster/errors"
)

const prefix = "/run/secrets/"

func GetSecret(name string) (string, error) {
	fmt.Printf("%#v\n", name) //nolint

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

		fmt.Printf("%#v\n", string(buff)) //nolint

		var decodedBuff []byte

		_, err = base64.StdEncoding.Decode(decodedBuff, buff)
		if err != nil {
			return "", errors.Wrap(err, "failed to decode secret")
		}

		fmt.Printf("%#v\n", string(decodedBuff)) //nolint

		return string(decodedBuff), nil
	}

	return name, nil
}
