package database

import (
	"app/database/schema"
	"fmt"
)

func Seed() {
	// Init Tables - Order Matters For Foreign Keys
	DB.MustExec(schema.User_Role)
	DB.MustExec(schema.User_Schema)
	// Seed Tables With Defaults
	SeedUserRoles()
	fmt.Println("\nSuccessfully Seeded Database")
}

func SeedUserRoles() {
	// Ignore Err On Duplicate
	_, err := DB.Exec(`INSERT INTO user_role (role, description) VALUES 
				('default', 'Assigned by default.'),
				('admin', 'Possess special privlidges.')
	`)
	if err != nil {
		fmt.Printf("Seed User Role: %v\n", err)
	}
}

func SeedUsers() {

}
