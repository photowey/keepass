package vault

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/photowey/keepass/configs"
)

const (
	fileMagic  = "KPAV"
	headerSize = 10
)

var (
	ErrInvalidFile        = errors.New("invalid vault file")
	ErrUnsupportedVersion = errors.New("unsupported vault format version")
	ErrDecryptFailed      = errors.New("unable to unlock vault")
)

type Codec interface {
	Version() uint16
	Encode(plaintext []byte, masterPassword string, cfg configs.Config) ([]byte, error)
	Decode(data []byte, masterPassword string) ([]byte, error)
}

var codecs = map[uint16]Codec{
	1: &codecV1{},
}

func Encode(formatVersion uint16, plaintext []byte, masterPassword string, cfg configs.Config) ([]byte, error) {
	codec, ok := codecs[formatVersion]
	if !ok {
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedVersion, formatVersion)
	}

	return codec.Encode(plaintext, masterPassword, cfg)
}

func Decode(data []byte, masterPassword string) ([]byte, error) {
	version, err := detectVersion(data)
	if err != nil {
		return nil, err
	}

	codec, ok := codecs[version]
	if !ok {
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedVersion, version)
	}

	return codec.Decode(data, masterPassword)
}

func detectVersion(data []byte) (uint16, error) {
	if len(data) < headerSize {
		return 0, ErrInvalidFile
	}

	if string(data[:4]) != fileMagic {
		return 0, ErrInvalidFile
	}

	return binary.BigEndian.Uint16(data[4:6]), nil
}
