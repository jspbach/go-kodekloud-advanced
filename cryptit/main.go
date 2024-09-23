package main

import (
	"fmt"
	"github.com/jspbach/golang-kodekloud-advanced/decrypt"
	"github.com/jspbach/golang-kodekloud-advanced/encrypt"
)

func main() {
	fmt.Println(encrypt.Nimbus("KodeKloud"))
	fmt.Println(decrypt.Nimbus(("KodeKloud")))
}
