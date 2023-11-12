/*
 * Copyright © 2023 the original author or authors.
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

package jsonz

import (
	"encoding/json"
	"io"
)

func String(body any) string {
	data, _ := StringE(body)

	return data
}

func StringE(body any) (string, error) {
	data, err := json.Marshal(body)

	return string(data), err
}

func Bytes(body any) []byte {
	data, _ := BytesE(body)

	return data
}

func BytesE(body any) ([]byte, error) {
	data, err := json.Marshal(body)

	return data, err
}

func Pretty(body any) string {
	data, _ := PrettyE(body)

	return data
}

func PrettyE(body any) (string, error) {
	bytes, err := json.MarshalIndent(body, "", "\t")

	return string(bytes), err
}

// ---------------------------------------------------------------- Decode

func DecodeStruct(reader io.Reader, target any) error {
	if err := json.NewDecoder(reader).Decode(target); err != nil {
		return err
	}

	return nil
}

func UnmarshalStructE(data []byte, structy any) error {
	if err := json.Unmarshal(data, structy); err != nil {
		return err
	}

	return nil
}

func UnmarshalStruct(data []byte, structy any) {
	_ = json.Unmarshal(data, structy)
}
