package main

import (
	"os"

	"github.com/yoshi-komoto/gitbucket-cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
