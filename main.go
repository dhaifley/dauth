// Package main is the entry point for the authentication services.
package main

import (
	"github.com/dhaifley/dauth/cmd"
	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
}
