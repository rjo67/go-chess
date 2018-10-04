package main

import (
	"fmt"

	"github.com/rjo67/chess/bitset"
)

func main() {
	bs := bitset.BitSet{Val: 4}
	fmt.Println(bs.ToString())
	fmt.Printf("bit 3 is set: %t\n", bs.IsSet(3))

}
