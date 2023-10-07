package userRoutes

import (
	"github.com/gofiber/fiber/v2"

	"app/api/auth"
	"app/models/user"
)

func SetUserRoutes(api fiber.Router) {
	userGroup := api.Group("/user")
	userGroup.Post("/login", user.Login) // Tracking Handled By Login Function
	userGroup.Post("/create", user.CreateUser)
	userGroup.Get("/get", auth.ValidateJWT, user.GetUser)
	userGroup.Put("/", auth.ValidateJWT, user.UpdateUser)
	userGroup.Delete("/", auth.ValidateJWT, auth.ValidateAdmin, user.DeleteUser)
	// userGroup.Put("/admin-update", auth.ValidateJWT, auth.ValidateAdmin, user.AdminUpdate)
	// userGroup.Get("/getall", auth.ValidateJWT, auth.ValidateAdmin, user.GetAll)

	userGroup.Get("/get-user-roles", auth.ValidateJWT, auth.ValidateAdmin, user.GetUserRoles)
}
