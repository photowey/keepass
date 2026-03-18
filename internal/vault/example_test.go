package vault

import (
	"fmt"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
)

func ExampleEncode() {
	env := home.Environment{
		RootDir:         "/tmp/example-home/.keepass",
		ConfigFile:      "/tmp/example-home/.keepass/keepass.config.json",
		DefaultVault:    "/tmp/example-home/.keepass/keepass.kp",
		ResolvedHomeDir: "/tmp/example-home",
	}

	cfg := configs.Default(env)
	cfg.Security.Argon2id.Time = 1
	cfg.Security.Argon2id.MemoryKiB = 8 * 1024
	cfg.Security.Argon2id.Threads = 1
	cfg.Security.Argon2id.KeyLength = 32

	data, err := Encode(1, []byte("secret payload"), "master-password", cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	plaintext, err := Decode(data, "master-password")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(plaintext))

	// Output:
	// secret payload
}
