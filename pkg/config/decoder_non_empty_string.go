package config

import (
	"fmt"
)

type nonEmptyStringDecoder string

func (nesd *nonEmptyStringDecoder) Decode(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("should not be empty")
	}

	*nesd = nonEmptyStringDecoder(value)
	return nil
}
