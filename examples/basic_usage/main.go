//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"go.g3deon.com/fieldmask"
)

type User struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Profile Profile `json:"profile"`
}

type Profile struct {
	Age     int    `json:"age"`
	Bio     string `json:"bio"`
	Website string `json:"website"`
}

func main() {
	// Create a user with complete data
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Profile: Profile{
			Age:     30,
			Bio:     "Software Developer",
			Website: "https://johndoe.com",
		},
	}

	// Create a field mask to keep only the name and profile.age
	mask := fieldmask.New("name", "profile.age")

	// Apply the mask - this will zero all other fields
	err := mask.Apply(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// The user now has only the specified fields
	fmt.Printf("%+v\n", user)
	// Output: {Name:John Doe Email: Profile:{Age:30 Bio: Website:}}
}
