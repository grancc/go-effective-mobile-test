package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gosubscription "github.com/grancc/go-effective-mobile-test"
	"github.com/grancc/go-effective-mobile-test/pkg/repository"
	"github.com/sirupsen/logrus"
)

func (h *Handler) createSubscription(c *gin.Context) {
	log := logrus.WithField("handler", "createSubscription")
	var input gosubscription.Subscription

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Subscriptions.CreateSubscription(input)

	if err != nil {
		log.WithError(err).Error("repository create failed")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.WithField("id", id).Info("subscription created")
	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getSubscriptionById(c *gin.Context) {
	log := logrus.WithField("handler", "getSubscriptionById")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	subscription, err := h.services.Subscriptions.GetSubscriptionById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("id", id).Warn("not found")
			NewErrorResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		log.WithError(err).Error("get by id failed")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.WithField("id", id).Info("subscription returned")
	c.JSON(http.StatusOK, subscription)
}

func (h *Handler) updateSubscriptionById(c *gin.Context) {
	log := logrus.WithField("handler", "updateSubscriptionById")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input gosubscription.Subscription
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subscription, err := h.services.Subscriptions.UpdateSubscriptionById(id, input)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("id", id).Warn("not found")
			NewErrorResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		log.WithError(err).Error("update failed")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.WithField("id", id).Info("subscription updated")
	c.JSON(http.StatusOK, subscription)
}

func (h *Handler) deleteSubscriptionById(c *gin.Context) {
	log := logrus.WithField("handler", "deleteSubscriptionById")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Subscriptions.DeleteSubscriptionById(id)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			log.WithField("id", id).Warn("not found")
			NewErrorResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		log.WithError(err).Error("delete failed")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.WithField("id", id).Info("subscription deleted")
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

func (h *Handler) listSubscriptions(c *gin.Context) {
	log := logrus.WithField("handler", "listSubscriptions")
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		NewErrorResponse(c, http.StatusBadRequest, "user_id query parameter is required")
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	subs, err := h.services.Subscriptions.ListSubscriptionsByUserId(userID)
	if err != nil {
		log.WithError(err).Error("list failed")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.WithField("count", len(subs)).WithField("user_id", userID).Info("list ok")
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) sumSubscriptions(c *gin.Context) {
	log := logrus.WithField("handler", "sumSubscriptions")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		NewErrorResponse(c, http.StatusBadRequest, "from and to query parameters are required (YYYY-MM-DD)")
		return
	}

	periodStart, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid from: use YYYY-MM-DD")
		return
	}
	periodEnd, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid to: use YYYY-MM-DD")
		return
	}
	if periodEnd.Before(periodStart) {
		NewErrorResponse(c, http.StatusBadRequest, "to must be on or after from")
		return
	}

	var userID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		u, err := uuid.Parse(userIDStr)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &u
	}

	var serviceName *string
	if name := c.Query("service_name"); name != "" {
		serviceName = &name
	}

	total, err := h.services.Subscriptions.SumSubscriptions(serviceName, userID, periodStart, periodEnd)
	if err != nil {
		log.WithError(err).Error("sum failed")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.WithFields(logrus.Fields{
		"from":  fromStr,
		"to":    toStr,
		"total": total,
	}).Info("sum ok")
	c.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
	})
}
