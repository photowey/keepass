package home

import (
	"fmt"
	"os"
)

func ExampleDetect() {
	originalValue, hadValue := os.LookupEnv(EnvKeepassHomePath)
	defer func() {
		if hadValue {
			_ = os.Setenv(EnvKeepassHomePath, originalValue)
			return
		}
		_ = os.Unsetenv(EnvKeepassHomePath)
	}()

	_ = os.Setenv(EnvKeepassHomePath, "/tmp/keepass-example")

	env, err := Detect()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(env.RootDir)
	fmt.Println(env.ConfigFile)
	fmt.Println(env.DefaultVault)

	// Output:
	// /tmp/keepass-example
	// /tmp/keepass-example/keepass.config.json
	// /tmp/keepass-example/keepass.kp
}
