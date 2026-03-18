package common_test

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
)

func ExampleParseTags() {
	fmt.Println(common.ParseTags(" code, ops , , prod "))

	// Output:
	// [code ops prod]
}
