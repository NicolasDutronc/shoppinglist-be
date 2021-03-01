package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/NicolasDutronc/shoppinglist-be/internal/user"
	"github.com/gin-gonic/gin"
)

// FindUserByIDHandler is a http handler for the FindByID service
func FindUserByIDHandler(srv user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := srv.FindByID(c.Request.Context(), id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// FindUserByNameHandler is a http handler for the FindByName service
func FindUserByNameHandler(srv user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		user, err := srv.FindByName(c.Request.Context(), name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// StoreUserHandler is a http handler for the Store service
func StoreUserHandler(srv user.Service) gin.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var req request

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user, err := srv.Store(c.Request.Context(), req.Name, req.Password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// UpdateUserNameHandler is a http handler for the UpdateName service
func UpdateUserNameHandler(srv user.Service) gin.HandlerFunc {
	type request struct {
		NewName string `json:"newName"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userID := c.Param("id")

		n, err := srv.UpdateName(c.Request.Context(), userID, req.NewName)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"number_of_updated": n,
		})
	}
}

// UpdateUserPasswordHandler is a http handler for the UpdatePassword service
func UpdateUserPasswordHandler(srv user.Service) gin.HandlerFunc {
	type request struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userID := c.Param("id")

		n, err := srv.UpdatePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"number_of_updated": n,
		})
	}
}

// DeleteUserHandler is a http handler for the Delete service
func DeleteUserHandler(srv user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		n, err := srv.Delete(c.Request.Context(), id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_deleted": n,
		})
	}
}

// AddPermissionsHandler is a http handler for the AddPermissions service
func AddPermissionsHandler(srv user.Service) gin.HandlerFunc {
	type request struct {
		Permissions []*user.Permission `json:"permissions"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userID := c.Param("id")

		n, err := srv.AddPermissions(c.Request.Context(), userID, req.Permissions...)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_updated": n,
		})
	}
}

// RemovePermissionsHandler is a http handler for the RemovePermissions service
func RemovePermissionsHandler(srv user.Service) gin.HandlerFunc {
	type request struct {
		Permissions []*user.Permission `json:"permissions"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userID := c.Param("id")

		n, err := srv.RemovePermissions(c.Request.Context(), userID, req.Permissions...)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_updated": n,
		})
	}
}

// LoginHandler is a http handler for the login service
func LoginHandler(srv user.Service) gin.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		User  *user.User `json:"user"`
		Token string     `json:"token"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user, token, err := srv.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.JSON(http.StatusOK, &response{
			User:  user,
			Token: token,
		})
	}
}

// AuthenticateMiddleware is a http middleware for the authenticate service
func AuthenticateMiddleware(srv user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Authorization header was not set"))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if !(len(headerParts) == 2 && headerParts[0] == "Bearer" && len(headerParts[1]) != 0) {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Authorization header was not set correctly"))
			return
		}

		token := headerParts[1]

		user, err := srv.Authenticate(c.Request.Context(), token)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}

// AuthorizationMiddleware is an authorization middleware to add in the handler chain of each handler with the correct permission configuration
func AuthorizationMiddleware(action string, resourceID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// rebuild action and resourceID based on dynamic params
		for _, param := range c.Params {
			if strings.Contains(resourceID, fmt.Sprintf(":%s", param.Key)) {
				resourceID = strings.Replace(resourceID, fmt.Sprintf(":%s", param.Key), param.Value, 1)
			}
			if strings.Contains(action, fmt.Sprintf(":%s", param.Key)) {
				action = strings.Replace(action, fmt.Sprintf(":%s", param.Key), param.Value, 1)
			}
		}

		// retrieve the current user
		currentUser, err := GetCurrentUser(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// check if the user has the permission
		if err := currentUser.Can(action, resourceID); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
