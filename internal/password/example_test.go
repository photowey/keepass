package password_test

import (
	"fmt"
	"strings"

	"github.com/photowey/keepass/internal/password"
)

func ExampleGenerate() {
	value, err := password.Generate(8, "ab")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(value))
	fmt.Println(strings.Trim(value, "ab") == "")

	// Output:
	// 8
	// true
}

func ExampleAlphabetForPreset() {
	alphabet, err := password.AlphabetForPreset(password.PresetCompatible)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(strings.Contains(alphabet, "-"))
	fmt.Println(strings.Contains(alphabet, "0"))

	// Output:
	// true
	// false
}
