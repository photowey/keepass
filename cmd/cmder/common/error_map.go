package common

import (
	"errors"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/vault"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, configs.ErrConfigNotFound),
		errors.Is(err, vault.ErrVaultNotInitialized):
		return WithExitCode(ExitCodeNotInitialized, err)
	case errors.Is(err, vault.ErrDecryptFailed):
		return WithExitCode(ExitCodeUnlockFailed, err)
	default:
		return err
	}
}
