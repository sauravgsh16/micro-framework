package main

import (
	"fmt"

	"github.com/sauravgsh16/micro-framework/node"
)

func main() {
	n1, _ := node.NewNode(1, 1)
	n2, _ := node.NewNode(2, 2)
	fmt.Printf("%+v\n", n1)
	fmt.Printf("%+v\n", n2)

	n1.AddReceiver(*n2, "broadcast")
}
