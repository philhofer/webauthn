package webauthn

import (
	"fmt"
	"time"
)

func errValidate(err error) error {
	return fmt.Errorf("error occurred validating the configuration: %w", err)
}

const (
	defaultTimeoutUVD = time.Millisecond * 120000
	defaultTimeout    = time.Millisecond * 300000
)
