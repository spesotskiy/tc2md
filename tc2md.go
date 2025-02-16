package main

import (
	"fmt"
	"tc2mdc"
)

func main() {
	fmt.Println("Convert test comments to a MD file.")
	mdText, err := tc2mdc.Convert(nil) //[]string{"// # GIVEN there is somethong"})
	fmt.Println(mdText, err)
}
