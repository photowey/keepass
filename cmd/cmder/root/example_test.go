package root_test

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/root"
)

func ExampleNewCommand() {
	cmd := root.NewCommand()

	fmt.Println(cmd.Use)
	fmt.Println(cmd.Short)
	fmt.Println(cmd.Version != "")

	// Output:
	// keepass
	// Secure local password manager
	// true
}
