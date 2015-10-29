package storage

import (
	"fmt"
	"strings"
)

// Storage factory
type StorageFactory func(arg map[string]string) (res Storage, err error)

func GuessStorage(arg string) (res Storage, err error) {
	split := strings.Split(arg, ":")
	opts := map[string]string{}

	if len(split) == 0 {
		err = fmt.Errorf("no storage configuration provided")
	}
	kind := split[0]

	if len(split) > 1 {
		split = strings.Split(split[1], ",")
		for _, v := range split {
			pair := strings.Split(v, "=")
			opts[pair[0]] = pair[1]
		}
	}

	switch kind {
	case "block":
		res, err = BlockStorageFactory(opts)
		return
	}

	return
}
