//go:build !windows

package cryptox

import (
	"fmt"
)

func WDecrypt(_ []byte) ([]byte, error) {
	return nil, fmt.Errorf("only available on windows")
}
