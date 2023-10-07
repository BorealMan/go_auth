package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"app/config"
)

var JWTSecretKey = []byte(config.JWT_SECRET)

func IssueJWT(userId int, userRole string) (string, error) {

	claims := jwt.MapClaims{
		"userId": fmt.Sprintf("%d", userId),
		"role":   userRole,
		"exp":    time.Now().Add(time.Minute * 3600).Unix(), // 1 Day
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	t, err := token.SignedString(JWTSecretKey)

	return t, err
}

func ValidateJWT(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	authHeader := headers["Authorization"]
	authToken := strings.Split(authHeader, " ")

	if len(authToken) != 2 {
		return c.SendStatus(403)
	}

	token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return c.SendStatus(403)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		// fmt.Printf("Claims:\nRole:%s\nUserId:%s\n", claims["role"], claims["userId"])
		if ok {
			c.Request().Header.Set("Role", fmt.Sprintf("%s", claims["role"]))
			c.Request().Header.Set("Userid", fmt.Sprintf("%s", claims["userId"]))
			// Add To Locals -- Plan To Convert All Methods To Use Locals Instead of Request Header
			c.Locals("userID", fmt.Sprintf("%s", claims["userId"]))
			c.Locals("role", fmt.Sprintf("%s", claims["role"]))
		}
		return c.Next()
	}

	return c.SendStatus(500)
}

// Use Validate JWT First
func ValidateAdmin(c *fiber.Ctx) error {
	// if c.GetReqHeaders()["Role"] != "admin" {
	// 	return c.SendStatus(403)
	// }
	if fmt.Sprintf("%s", c.Locals("role")) != "admin" {
		return c.SendStatus(403)
	}
	return c.Next()
}
