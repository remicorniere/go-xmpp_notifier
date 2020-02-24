package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello gophers lel")
	fmt.Println(len(os.Args), os.Args)
}
