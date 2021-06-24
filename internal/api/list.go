package api

import (
	"net/http"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	"github.com/gin-gonic/gin"
)

// FindListByIDHandler retrieves a list based on the id passed in params
func FindListByIDHandler(srv list.FinderByID) gin.HandlerFunc {
	type response struct {
		ID        string       `json:"id"`
		CreatedAt time.Time    `json:"created_at"`
		UpdatedAt time.Time    `json:"updated_at"`
		Name      string       `json:"name"`
		Items     []*list.Item `json:"items"`
	}

	return func(c *gin.Context) {
		id := c.Param("id")
		list, err := srv.FindListByID(c.Request.Context(), id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"list": &response{
				ID:        list.ID.Hex(),
				CreatedAt: list.CreatedAt,
				UpdatedAt: list.UpdatedAt,
				Name:      list.Name,
				Items:     list.Items,
			},
		})
	}
}

// FindAllListsHandler returns all lists
func FindAllListsHandler(srv list.Finder) gin.HandlerFunc {
	type encodedList struct {
		ID        string       `json:"id"`
		CreatedAt time.Time    `json:"created_at"`
		UpdatedAt time.Time    `json:"updated_at"`
		Name      string       `json:"name"`
		Items     []*list.Item `json:"items"`
	}

	type response struct {
		Lists []*encodedList `json:"lists"`
	}

	return func(c *gin.Context) {
		lists, err := srv.FindAllLists(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// get current user
		currentUser, err := GetCurrentUser(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		response := &response{
			Lists: []*encodedList{},
		}

		for _, list := range lists {
			// if the user has the permission, the list is appended
			if err := currentUser.Can("read", "list-"+list.ID.Hex()); err == nil {
				response.Lists = append(response.Lists, &encodedList{
					ID:        list.ID.Hex(),
					CreatedAt: list.CreatedAt,
					UpdatedAt: list.UpdatedAt,
					Name:      list.Name,
					Items:     list.Items,
				})
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetInventoryHandler(srv list.Service) gin.HandlerFunc {
	type summary struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Length    int       `json:"length"`
	}

	type response struct {
		Lists []*summary `json:"inventory"`
	}
	return func(c *gin.Context) {
		lists, err := srv.FindAllLists(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// get current user
		currentUser, err := GetCurrentUser(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		response := &response{
			Lists: []*summary{},
		}

		for _, list := range lists {
			// if the user has the permission, the list is appended
			if err := currentUser.Can("read", "list-"+list.ID.Hex()); err == nil {
				response.Lists = append(response.Lists, &summary{
					ID:        list.ID.Hex(),
					CreatedAt: list.CreatedAt,
					UpdatedAt: list.UpdatedAt,
					Name:      list.Name,
					Length:    len(list.Items),
				})
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// DeleteListHandler removes a list based on its id
func DeleteListHandler(srv list.Deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		n, err := srv.DeleteList(c.Request.Context(), id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_deleted": n,
		})
	}
}

// StoreListHandler creates a new list and returns it
func StoreListHandler(srv list.Creator) gin.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	type response struct {
		ID        string       `json:"id"`
		CreatedAt time.Time    `json:"created_at"`
		UpdatedAt time.Time    `json:"updated_at"`
		Name      string       `json:"name"`
		Items     []*list.Item `json:"items"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		list, err := srv.StoreList(c.Request.Context(), req.Name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"list": &response{
				ID:        list.ID.Hex(),
				CreatedAt: list.CreatedAt,
				UpdatedAt: list.UpdatedAt,
				Name:      list.Name,
				Items:     list.Items,
			},
		})
	}
}

// RemoveAllItemsHandler removes all items from the list but does not delete it
func RemoveAllItemsHandler(srv list.Clearer) gin.HandlerFunc {
	return func(c *gin.Context) {
		listID := c.Param("id")

		n, err := srv.RemoveAllItems(c.Request.Context(), listID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_deleted": n,
		})
	}
}

// AddItemHandler returns a handler for adding an item
func AddItemHandler(srv list.ItemAdder) gin.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
	}

	return func(c *gin.Context) {
		listID := c.Param("id")
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// store new item
		item, err := srv.AddItem(c.Request.Context(), listID, req.Name, req.Quantity)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"item": &item,
		})
	}
}

// UpdateItemHandler returns a handler for updating an item
func UpdateItemHandler(srv list.ItemUpdater) gin.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Quantity    string `json:"quantity"`
		NewName     string `json:"new_name"`
		NewQuantity string `json:"new_quantity"`
	}
	return func(c *gin.Context) {
		listID := c.Param("id")
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		n, err := srv.UpdateItem(c.Request.Context(), listID, req.Name, req.Quantity, req.NewName, req.NewQuantity)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_updated": n,
		})
	}
}

// ToggleItemHandler returns a handler for toggling an item
func ToggleItemHandler(srv list.ItemToggler) gin.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
		Value    bool   `json:"value"`
	}
	return func(c *gin.Context) {
		listID := c.Param("id")
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		n, err := srv.ToggleItem(c.Request.Context(), listID, req.Name, req.Quantity, req.Value)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_toggled": n,
		})
	}
}

// RemoveItemHandler returns a handler for removing an item
func RemoveItemHandler(srv list.ItemRemover) gin.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
	}
	return func(c *gin.Context) {
		listID := c.Param("id")
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		n, err := srv.RemoveItem(c.Request.Context(), listID, req.Name, req.Quantity)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_deleted": n,
		})
	}
}
