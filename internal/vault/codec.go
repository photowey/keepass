/*
 * Copyright © 2023-present the keepass authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
