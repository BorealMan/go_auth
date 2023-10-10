package user

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"app/api/auth"
	"app/database"
)

var DEBUG = true

type User struct {
	ID              int            `db:"id" json:"id"`
	Email           string         `db:"email" json:"email" validate:"required,email"`
	Password        string         `db:"password" json:"password" validate:"omitempty,min=1,max=100"`
	RoleID          uint           `db:"role_id" json:"role_id" validate:"omitempty,number"`
	Role            User_Role      `json:"user_role" validate:"omitempty"`
	Phone           sql.NullString `db:"phone" json:"phone" validate:"omitempty,e164"`
	Account_Enabled *bool          `db:"account_enabled" json:"account_enabled" validate:"omitempty"`
	Created         time.Time      `db:"created" json:"created"`
	Updated_at      time.Time      `db:"updated_at" json:"updated_at"`
}

type User_Role struct {
	ID          int       `json:"id"`
	Role        string    `json:"role"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated_at  time.Time `json:"updated_at"`
}

// Test Type
type User_Role_Join struct {
	User_ID         int            `db:"id" json:"id"`
	Email           string         `db:"email" json:"email"`
	Phone           sql.NullString `db:"phone" json:"phone"`
	Account_enabled *bool          `db:"account_enabled"`
	// Role Values
	RoleID      uint64 `db:"role_id" json:"role_id"`
	Role        string `db:"role" json:"role"`
	Description string `db:"description" json:"description"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1,max=100"`
}

func Login(c *fiber.Ctx) error {

	r := new(LoginRequest)
	err := c.BodyParser(r)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": "Invalid Input Fields"})
	}

	if DEBUG {
		log.Printf("Login -> Username: %s\n", r.Email)
	}

	// Process Input
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Password = strings.TrimSpace(r.Password)

	// Validation
	err = Validate(r)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed Validation"})
	}

	// Look Up User by Username
	user := new(User)
	err = database.DB.Get(user, `SELECT * from user WHERE email = ?`, r.Email)

	if err != nil {
		if DEBUG {
			log.Printf(`Login Error: %v`, err)
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Username or Password"})
	}

	// Verify User Account is Enabled If Found
	if !*user.Account_Enabled {
		return c.Status(403).JSON(fiber.Map{"error": "User Account Is Disabled"})
	}

	// Check Password
	pld := r.Email + r.Password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pld)) != nil {
		log.Println("Login Bcrypt Error")
		return c.Status(500).JSON(fiber.Map{"error": "Invalid Username or Password"})
	}

	// Get Role
	role := new(User_Role)
	err = database.DB.Get(role, `SELECT * from user_role WHERE id = ?`, user.RoleID)

	if err != nil {
		log.Printf(`Login Failed To Retreive Role: %v`, err)
		return c.Status(500).JSON(fiber.Map{"error": "Unknown Issue - Please Try Again Later"})
	}

	user.Role = *role

	// Create Token
	t, err := auth.IssueJWT(user.ID, user.Role.Role)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed - Try To Login"})
	}

	user.Password = ""

	return c.Status(200).JSON(fiber.Map{"user": user, "token": t})
}

type CreateUserRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func CreateUser(c *fiber.Ctx) error {
	// Parse Request
	r := new(CreateUserRequest)
	err := c.BodyParser(r)
	if err != nil {
		if DEBUG {
			log.Printf(`Create User Failed To Parse: %v`, err)
		}
		return c.Status(400).JSON(fiber.Map{"error": err})
	}
	// Process Data
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Password = strings.TrimSpace(r.Password)

	// Validate
	err = Validate(r)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	user := new(User)
	// Look Up If Username Exists
	err = database.DB.Get(user, `SELECT id from user WHERE email = ?`, r.Email)

	if user.ID != 0 {
		log.Println(`Create User Username Already Exists`)
		return c.Status(400).JSON(fiber.Map{"error": "Username Already Exists"})
	}

	// Hash Password
	pld := r.Email + r.Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(pld), 7)
	if err != nil {
		log.Printf(`Create User Error Hashing Password: %v`, err)
		return c.Status(500).JSON(fiber.Map{"error": "Unable To Create User At This Time"})
	}

	user = new(User)
	user.Email = r.Email
	user.Password = string(bytes)
	// Set Default Role ID
	user.RoleID = 1
	// Create User
	_, err = database.DB.NamedExec(`INSERT INTO user (email, password, role_id) 
									VALUES (:email, :password, :role_id)`, user)
	if err != nil {
		log.Printf(`Create User Database Error: %v`, err)
		return c.Status(500).JSON(fiber.Map{"error": "Unable To Create User At This Time"})
	}

	// Build User Return Object
	err = database.DB.Get(user, `SELECT * from user WHERE email = ?`, user.Email)

	if err != nil {
		if DEBUG {
			log.Println(err)
		}
		return c.Status(500).JSON(fiber.Map{"error": "Unknown Error Try To Login"})
	}

	role := new(User_Role)
	database.DB.Get(role, `SELECT * from user_role WHERE id = ?`, user.RoleID)

	// Assign Role
	user.Role = *role

	if DEBUG {
		log.Println(user)
		log.Println(role)
	}

	// Create Token
	t, err := auth.IssueJWT(user.ID, user.Role.Role)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed - Try To Login"})
	}

	user.Password = ""

	return c.Status(201).JSON(fiber.Map{"token": t, "user": user})
}

func GetUser(c *fiber.Ctx) error {
	user_id, _ := strconv.ParseUint(fmt.Sprintf("%s", c.Locals("userID")), 10, 64)
	user := new(User)
	err := database.DB.Get(user, `SELECT * FROM user WHERE id = ?`, user_id)
	if err != nil {
		if DEBUG {
			log.Println("Get User: Unable To Find User")
		}
		return c.Status(400).JSON(fiber.Map{"error": "Unable To Find User"})
	}
	role := new(User_Role)
	err = database.DB.Get(role, `SELECT * FROM user_role WHERE id = ?`, user.RoleID)
	if err != nil {
		if DEBUG {
			log.Println("Get User: Unable To Find Role")
		}
		return c.Status(500).JSON(fiber.Map{"error": "Unknown Error"})
	}
	user.Role = *role
	user.Password = ""
	return c.Status(200).JSON(fiber.Map{"user": user})
}

func UpdateUser(c *fiber.Ctx) error {
	return nil
}

func DeleteUser(c *fiber.Ctx) error {
	return nil
}

func GetUserRoles(c *fiber.Ctx) error {
	var userRoles []User_Role
	err := database.DB.Select(userRoles, `SELECT * FROM user_role`)
	if err != nil {
		if DEBUG {
			log.Printf("Get User Roles: %v", err)
		}
		return c.Status(500).JSON(fiber.Map{"error": "Unable To Fetch User Roles"})
	}
	return c.Status(200).JSON(fiber.Map{"user_roles": userRoles})
}

func Validate(r interface{}) error {
	// Validation
	validate := validator.New()
	validate_err := validate.Struct(r)

	if validate_err != nil {
		if DEBUG {
			log.Println("Validation Error: ", validate_err)
		}
		return validate_err
	}
	return nil
}
