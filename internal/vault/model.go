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

import "time"

type Document struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Entries   []Entry   `json:"entries"`
}

type Entry struct {
	Alias             string    `json:"alias"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	URI               string    `json:"uri,omitempty"`
	Note              string    `json:"note,omitempty"`
	Tags              []string  `json:"tags,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
}

func NewDocument(now time.Time) *Document {
	return &Document{
		Version:   1,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
		Entries:   []Entry{},
	}
}
