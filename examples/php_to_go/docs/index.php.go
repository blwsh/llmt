//go:build exclude

package main

import (
	"fmt"
	"log"

	"benwatson/myphpproject/auth"
	"greeter"
)

func main() {
	authorizer := auth.NewAuthorizer("ben")

	if !authorizer.Authorize("ben") {
		log.Fatal("You are authorized!")
	}

	greeter := greeter.NewGreeter()

	fmt.Println(greeter.Greet())
}
