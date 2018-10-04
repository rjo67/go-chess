package main

import (
	"fmt"

	"github.com/rjo67/chess/position"
)

func main() {
	posn := position.StartPosition()

	fmt.Printf("Start position:\n%s", posn.ToString())

}
