package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func Generate(length int, alphabet string) (string, error) {
	if length <= 0 {
		return "", errors.New("password length must be positive")
	}

	alphabet = strings.TrimSpace(alphabet)
	if alphabet == "" {
		return "", errors.New("alphabet cannot be blank")
	}

	limit := big.NewInt(int64(len(alphabet)))
	var builder strings.Builder
	builder.Grow(length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", fmt.Errorf("generate password: %w", err)
		}

		builder.WriteByte(alphabet[index.Int64()])
	}

	return builder.String(), nil
}
