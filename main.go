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

package main

import (
	App "github.com/photowey/keepass/cmd/cmder/root"
)

// $ go install github.com/go-delve/delve/cmd/dlv@v1.20.0

func main() {
	// $ keepass add -n photowey -a github -u photowey@github.com -p hello.github -i https://github.com -t "the username and password of github.com website"
	// $ keepass add photowey -a github -u photowey@github.com -p hello.github -i https://github.com -t "the username and password of github.com website"
	// $ keepass add < testdata.json
	App.Run()
}
