package api

import (
	"fmt"

	"github.com/NicolasDutronc/shoppinglist-be/internal/user"
	"github.com/gin-gonic/gin"
)

// GetCurrentUser retrieves the current user from the gin context where it was set by the authentication middleware
func GetCurrentUser(c *gin.Context) (*user.User, error) {
	interfaceUser, exists := c.Get("currentUser")
	if !exists {
		return nil, fmt.Errorf("current user not found in the gin context")
	}

	user, ok := interfaceUser.(*user.User)
	if !ok {
		return nil, fmt.Errorf("current user was not compatible with the user struct")
	}

	return user, nil
}
