package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	PresetCompatible        = "compatible"
	PresetSymbols           = "symbols"
	PresetStrictHighEntropy = "strict-high-entropy"
)

var presetAlphabets = map[string]string{
	PresetCompatible:        "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789-_",
	PresetSymbols:           "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789-_!@#$%^&*+=",
	PresetStrictHighEntropy: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[]{}:,.?",
}

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

func AlphabetForPreset(preset string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(preset))
	if normalized == "" {
		return "", errors.New("preset cannot be blank")
	}

	alphabet, ok := presetAlphabets[normalized]
	if !ok {
		return "", fmt.Errorf("unsupported password preset %q", preset)
	}

	return alphabet, nil
}
