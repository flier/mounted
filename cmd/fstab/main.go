package main

import (
	"fmt"
	"os"

	"github.com/flier/mounted"
)

func main() {
	fstab, err := mounted.FileSystems()

	if err != nil {
		fmt.Printf("fail to get mounted file systems, %s", err)

		os.Exit(-1)
	}

	for _, fs := range fstab {
		fmt.Printf("%s\n", fs)
	}
}
