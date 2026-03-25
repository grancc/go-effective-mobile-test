package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gosubscription "github.com/grancc/go-effective-mobile-test"
)

func (h *Handler) createSubscription(c *gin.Context) {
	var input gosubscription.Subscription

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Subscriptions.CreateSubscription(input)

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
