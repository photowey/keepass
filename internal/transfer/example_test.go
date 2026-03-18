package transfer

import (
	"fmt"

	"github.com/photowey/keepass/internal/vault"
)

func ExampleNormalizeConflictStrategy() {
	strategy, err := NormalizeConflictStrategy(" OVERWRITE ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(strategy)

	// Output:
	// overwrite
}

func ExampleSortEntries() {
	entries := []vault.Entry{{Alias: "zeta"}, {Alias: "alpha"}, {Alias: "kappa"}}
	SortEntries(entries)

	for _, entry := range entries {
		fmt.Println(entry.Alias)
	}

	// Output:
	// alpha
	// kappa
	// zeta
}
