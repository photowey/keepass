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
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/photowey/keepass/configs"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

type codecV1 struct{}

type v1Header struct {
	KDF      string                 `json:"kdf"`
	Cipher   string                 `json:"cipher"`
	Salt     string                 `json:"salt"`
	Nonce    string                 `json:"nonce"`
	Argon2id configs.Argon2idConfig `json:"argon2id"`
}

func (c *codecV1) Version() uint16 {
	return 1
}

func (c *codecV1) Encode(plaintext []byte, masterPassword string, cfg configs.Config) ([]byte, error) {
	if masterPassword == "" {
		return nil, errors.New("master password cannot be blank")
	}

	salt, err := randomBytes(16)
	if err != nil {
		return nil, err
	}

	nonce, err := randomBytes(chacha20poly1305.NonceSizeX)
	if err != nil {
		return nil, err
	}

	header := v1Header{
		KDF:      "argon2id",
		Cipher:   "xchacha20poly1305",
		Salt:     base64.StdEncoding.EncodeToString(salt),
		Nonce:    base64.StdEncoding.EncodeToString(nonce),
		Argon2id: cfg.Security.Argon2id,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("marshal header: %w", err)
	}

	key := deriveKey([]byte(masterPassword), salt, cfg.Security.Argon2id)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, headerBytes)

	data := make([]byte, headerSize+len(headerBytes)+len(ciphertext))
	copy(data[:4], []byte(fileMagic))
	binary.BigEndian.PutUint16(data[4:6], c.Version())
	binary.BigEndian.PutUint32(data[6:10], uint32(len(headerBytes)))
	copy(data[10:10+len(headerBytes)], headerBytes)
	copy(data[10+len(headerBytes):], ciphertext)

	return data, nil
}

func (c *codecV1) Decode(data []byte, masterPassword string) ([]byte, error) {
	if masterPassword == "" {
		return nil, errors.New("master password cannot be blank")
	}

	if len(data) < headerSize {
		return nil, ErrInvalidFile
	}

	headerLength := int(binary.BigEndian.Uint32(data[6:10]))
	if headerLength <= 0 || len(data) < headerSize+headerLength {
		return nil, ErrInvalidFile
	}

	headerBytes := data[10 : 10+headerLength]
	ciphertext := data[10+headerLength:]
	if len(ciphertext) == 0 {
		return nil, ErrInvalidFile
	}

	var header v1Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}

	if header.KDF != "argon2id" || header.Cipher != "xchacha20poly1305" {
		return nil, ErrInvalidFile
	}

	if header.Argon2id.Time == 0 || header.Argon2id.MemoryKiB == 0 || header.Argon2id.Threads == 0 || header.Argon2id.KeyLength == 0 {
		return nil, ErrInvalidFile
	}

	salt, err := base64.StdEncoding.DecodeString(header.Salt)
	if err != nil {
		return nil, fmt.Errorf("decode salt: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(header.Nonce)
	if err != nil {
		return nil, fmt.Errorf("decode nonce: %w", err)
	}

	if len(salt) == 0 || len(nonce) != chacha20poly1305.NonceSizeX {
		return nil, ErrInvalidFile
	}

	key := deriveKey([]byte(masterPassword), salt, header.Argon2id)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, headerBytes)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

func deriveKey(masterPassword, salt []byte, cfg configs.Argon2idConfig) []byte {
	return argon2.IDKey(masterPassword, salt, cfg.Time, cfg.MemoryKiB, cfg.Threads, cfg.KeyLength)
}

func randomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("read secure random bytes: %w", err)
	}

	return buf, nil
}
