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
	users.PUT("/:id/name", UpdateUserNameHandler(userSrv))
	users.PUT("/:id/password", UpdateUserPasswordHandler(userSrv))
	users.DELETE("/:id", DeleteUserHandler(userSrv))
	users.PUT("/:id/permissions/add", AddPermissionsHandler(userSrv))
	users.PUT("/:id/permissions/remove", RemovePermissionsHandler(userSrv))

	lists := restricted.Group("/lists")
	lists.GET("", FindAllListsHandler(listSrv))
	lists.POST("", StoreListHandler(listSrv))
	lists.GET("/inventory", GetInventoryHandler(listSrv))

	listI := lists.Group("/:id")
	listI.GET("", AuthorizationMiddleware("read", "list-:id"), FindListByIDHandler(listSrv))
	listI.POST("", AuthorizationMiddleware("write", "list-:id"), AddItemHandler(listSrv))
	listI.PUT("", AuthorizationMiddleware("write", "list-:id"), UpdateItemHandler(listSrv))
	listI.PUT("/toggle", AuthorizationMiddleware("write", "list-:id"), ToggleItemHandler(listSrv))
	listI.PUT("/delete", AuthorizationMiddleware("write", "list-:id"), RemoveItemHandler(listSrv))
	listI.PUT("/clear", AuthorizationMiddleware("write", "list-:id"), RemoveAllItemsHandler(listSrv))
	listI.DELETE("", AuthorizationMiddleware("write", "list-:id"), DeleteListHandler(listSrv))

	hubGroup := restricted.Group("/hub")
	hubGroup.GET("/connect", hub.SubscribeJSONHandler(h))
	hubGroup.POST("/subscribe", hub.SubscriptionHandler(h))
	hubGroup.POST("/unsubscribe", hub.UnsubscriptionHandler(h))

	return r
}
