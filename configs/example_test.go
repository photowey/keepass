package configs_test

import (
	"fmt"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
)

func ExampleConfig_ResolveVaultPath() {
	env := home.Environment{
		RootDir:         "/tmp/example-home/.keepass",
		ConfigFile:      "/tmp/example-home/.keepass/keepass.config.json",
		DefaultVault:    "/tmp/example-home/.keepass/keepass.kp",
		ResolvedHomeDir: "/tmp/example-home",
	}

	cfg := configs.Default(env)
	fmt.Println(cfg.ResolveVaultPath(env))

	// Output:
	// /tmp/example-home/.keepass/keepass.kp
}

func ExamplePasswordGenerator_EffectiveAlphabet() {
	generator := configs.PasswordGenerator{Alphabet: "abc123"}
	alphabet, err := generator.EffectiveAlphabet()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(alphabet)

	// Output:
	// abc123
}
