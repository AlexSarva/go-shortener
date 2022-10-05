package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("hello World")
	os.Exit(1) // want "found os.Exit method in main function"
}
