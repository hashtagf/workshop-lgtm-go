package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

type UserHandler struct {
	Ch *amqp091.Channel
}

func (h *UserHandler) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Publish the user creation event to a message queue (rabbitmq)
		PublishUserCreationEvent(User{
			Name:  req.Name,
			Email: req.Email,
		}, c.Request.Context(), h.Ch)

		c.JSON(http.StatusCreated, CreateUserResponse{
			Message: "User created successfully",
			User: User{
				ID:    1, // Mock ID
				Name:  req.Name,
				Email: req.Email,
			},
		})
	}
}
