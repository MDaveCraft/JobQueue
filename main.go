package main

import (
	nanoid "github.com/mdavecraft/job-queue/nanoid"
)

func main() {
	// Example usage of nanoid
	id, err := nanoid.New(21)
	if err != nil {
		panic(err)
	}
	println("Generated ID:", id)
}