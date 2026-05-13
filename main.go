package main

import (
	"os"

	"github.com/yoshi-komoto/gitbucket-cli/cmd/gb"
)

func main() {
	os.Exit(gb.Execute())
}
