package auth

import (
	"net/http"

	"order-management/config"

	"github.com/gin-gonic/gin"
)

// Controller handles authentication HTTP requests.
type Controller struct {
	Config config.Config
}

// NewController creates an authentication controller.
func NewController(cfg config.Config) *Controller {
	return &Controller{
		Config: cfg,
	}
}

// Login authenticates a user and returns a JWT.
func (c *Controller) Login(ctx *gin.Context) {
	var request LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username and password are required",
		})
		return
	}

	user, err := Authenticate(
		c.Config,
		request.Username,
		request.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	token, err := GenerateToken(c.Config, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate authentication token",
		})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Type:  "Bearer",
		Role:  user.Role,
	})
}
