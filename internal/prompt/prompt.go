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

package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type Prompter struct {
	reader     *bufio.Reader
	out        io.Writer
	file       *os.File
	isTerminal bool
}

func New(in io.Reader, out io.Writer) *Prompter {
	p := &Prompter{
		reader: bufio.NewReader(in),
		out:    out,
	}

	if file, ok := in.(*os.File); ok {
		p.file = file
		p.isTerminal = term.IsTerminal(int(file.Fd()))
	}

	return p
}

func (p *Prompter) Ask(label string) (string, error) {
	if _, err := fmt.Fprintf(p.out, "%s: ", label); err != nil {
		return "", err
	}

	line, err := p.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func (p *Prompter) AskOptional(label string) (string, error) {
	return p.Ask(label)
}

func (p *Prompter) AskDefault(label, value string) (string, error) {
	if _, err := fmt.Fprintf(p.out, "%s [%s]: ", label, value); err != nil {
		return "", err
	}

	line, err := p.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return value, nil
	}

	return trimmed, nil
}

func (p *Prompter) AskSecret(label string) (string, error) {
	if _, err := fmt.Fprintf(p.out, "%s: ", label); err != nil {
		return "", err
	}

	if p.isTerminal && p.file != nil {
		value, err := term.ReadPassword(int(p.file.Fd()))
		if _, newlineErr := fmt.Fprintln(p.out); newlineErr != nil && err == nil {
			err = newlineErr
		}
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(value)), nil
	}

	line, err := p.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func (p *Prompter) AskSecretWithConfirmation(label, confirmationLabel string) (string, error) {
	first, err := p.AskSecret(label)
	if err != nil {
		return "", err
	}

	second, err := p.AskSecret(confirmationLabel)
	if err != nil {
		return "", err
	}

	if first != second {
		return "", fmt.Errorf("%s does not match confirmation", label)
	}

	return first, nil
}

func (p *Prompter) Confirm(label string, defaultYes bool) (bool, error) {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}

	answer, err := p.Ask(fmt.Sprintf("%s %s", label, suffix))
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer == "" {
		return defaultYes, nil
	}

	return answer == "y" || answer == "yes", nil
}
