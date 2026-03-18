package files_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/photowey/keepass/pkg/files"
)

func ExampleWriteFileAtomic() {
	dir, err := os.MkdirTemp("", "keepass-files-example-")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "vault.txt")
	if err := files.WriteFileAtomic(path, []byte("vault-data"), 0o600); err != nil {
		fmt.Println(err)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(files.Exists(path))
	fmt.Println(string(data))

	// Output:
	// true
	// vault-data
}
