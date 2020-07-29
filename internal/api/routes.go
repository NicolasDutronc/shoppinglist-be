package api

import (
	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	"github.com/NicolasDutronc/shoppinglist-be/internal/user"
	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers the routes to the router
func SetupRoutes(userSrv user.Service, listSrv list.Service, h hub.Hub) *gin.Engine {
	r := gin.Default()

	r.POST("/api/v1/login", LoginHandler(userSrv))

	restricted := r.Group("/api/v1")
	restricted.Use(AuthenticateMiddleware(userSrv))

	users := restricted.Group("/users")
	users.GET("/id/:id", FindUserByIDHandler(userSrv))
	users.GET("/name/:name", FindUserByNameHandler(userSrv))
	users.POST("", StoreUserHandler(userSrv))
	users.PUT("/name", UpdateUserNameHandler(userSrv))
	users.PUT("/password", UpdateUserPasswordHandler(userSrv))
	users.DELETE("/:id", DeleteUserHandler(userSrv))

	lists := restricted.Group("/lists")
	lists.GET("", FindAllListsHandler(listSrv))
	lists.POST("", StoreListHandler(listSrv))

	listI := lists.Group("/:id")
	listI.GET("", FindListByIDHandler(listSrv))
	listI.POST("", AddItemHandler(listSrv))
	listI.PUT("", UpdateItemHandler(listSrv))
	listI.PUT("/toggle", ToggleItemHandler(listSrv))
	listI.PUT("/delete", RemoveItemHandler(listSrv))
	listI.PUT("/clear", RemoveAllItemsHandler(listSrv))
	listI.DELETE("", DeleteListHandler(listSrv))

	hubGroup := restricted.Group("/hub")
	hubGroup.GET("/connect", hub.SubscribeJSONHandler(h))
	hubGroup.POST("/subscribe", hub.SubscriptionHandler(h))
	hubGroup.POST("/unsubscribe", hub.UnsubscriptionHandler(h))

	return r
}
