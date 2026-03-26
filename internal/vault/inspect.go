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
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/photowey/keepass/configs"
	"golang.org/x/crypto/chacha20poly1305"
)

type Metadata struct {
	FormatVersion uint16                  `json:"format_version"`
	KDF           string                  `json:"kdf,omitempty"`
	Cipher        string                  `json:"cipher,omitempty"`
	Argon2id      *configs.Argon2idConfig `json:"argon2id,omitempty"`
}

func InspectFile(path string) (Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Metadata{}, ErrVaultNotInitialized
		}

		return Metadata{}, fmt.Errorf("read vault: %w", err)
	}

	version, err := detectVersion(data)
	if err != nil {
		return Metadata{}, err
	}

	switch version {
	case 1:
		return inspectV1(data)
	default:
		return Metadata{}, fmt.Errorf("%w: %d", ErrUnsupportedVersion, version)
	}
}

func inspectV1(data []byte) (Metadata, error) {
	if len(data) < headerSize {
		return Metadata{}, ErrInvalidFile
	}

	headerLength := int(binary.BigEndian.Uint32(data[6:10]))
	if headerLength <= 0 || len(data) < headerSize+headerLength {
		return Metadata{}, ErrInvalidFile
	}

	headerBytes := data[10 : 10+headerLength]

	var header v1Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Metadata{}, fmt.Errorf("decode header: %w", err)
	}

	if header.KDF != "argon2id" || header.Cipher != "xchacha20poly1305" {
		return Metadata{}, ErrInvalidFile
	}

	salt, err := base64.StdEncoding.DecodeString(header.Salt)
	if err != nil || len(salt) == 0 {
		return Metadata{}, ErrInvalidFile
	}

	nonce, err := base64.StdEncoding.DecodeString(header.Nonce)
	if err != nil || len(nonce) != chacha20poly1305.NonceSizeX {
		return Metadata{}, ErrInvalidFile
	}

	if header.Argon2id.Time == 0 || header.Argon2id.MemoryKiB == 0 || header.Argon2id.Threads == 0 || header.Argon2id.KeyLength == 0 {
		return Metadata{}, ErrInvalidFile
	}

	argon2id := header.Argon2id
	return Metadata{
		FormatVersion: 1,
		KDF:           header.KDF,
		Cipher:        header.Cipher,
		Argon2id:      &argon2id,
	}, nil
}
