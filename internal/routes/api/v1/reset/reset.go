package reset

import (
	"context"
	"net/http"
	"time"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetResets(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	defer cancel()

	var req models.GetResetsInput

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error binding req")
		c.Status(http.StatusBadRequest)
		return
	}

	resets, err := database.GetResetsFromCounterId(ctx, req.CounterID, &req.Filter)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting resets for counter id")
		c.Status(http.StatusBadRequest)
		return
	}

	finalResets := make([]*models.GetResetsResetPart, len(resets))

	isCounterOwner := req.CounterID == user.Uuid

	for index, reset := range resets {
		isResetOwner := reset.ResetID == user.Uuid

		finalResets[index] = &models.GetResetsResetPart{
			Id:          reset.ResetID,
			Username:    reset.ResetUsername,
			Description: reset.ResetDescription,
			Timestamp:   reset.ResetOccurredAt,
			CanDelete:   (users.IsEboard(user) || users.IsActiveRTP(user) || isCounterOwner || isResetOwner),
		}
	}

	var hasMore bool = false
	var timestamp int64 = 0

	if len(finalResets) > 1 {
		hasMore = (len(finalResets) > (database.MAX_RESETS_IN_QUERY)-1)
		timestamp = finalResets[len(finalResets)-1].Timestamp
	}

	c.JSON(http.StatusOK, models.GetResetsOutput{
		Resets:    finalResets,
		HasMore:   hasMore,
		Timestamp: &timestamp,
	})
}

func CreateReset(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	defer cancel()

	var req models.CreateResetInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("failed to cast error %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if req.CounterID == "01a119b8-319a-7db2-9585-866bd641ac8c" {
		logging.Logger.Info("You cannot stop this")
		c.Status(http.StatusBadRequest)
		return
	}

	if _, err := database.CreateReset(ctx, user.Uuid, user.Username, &req); err != nil {
		logging.Logger.Warnf("failed to create error %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if _, err := database.UpdateCounterLastReset(ctx, req.CounterID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error updating last reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteReset(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting CSHAuth reset")
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	defer cancel()

	var req models.DeleteResetInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error compiling reset")
		c.Status(http.StatusBadRequest)
		return
	}

	counterID, err := database.GetCounterFromReset(ctx, req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting counter from reset")
		c.Status(http.StatusBadRequest)
		return
	}

	counterOwner, err := database.GetCounterOwner(ctx, counterID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting counter owner")
		c.Status(http.StatusBadRequest)
		return
	}

	resetOwner, err := database.GetResetOwner(ctx, req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error deleting reset")
		c.Status(http.StatusBadRequest)
		return
	}

	isCounterOwner := counterOwner == user.Uuid
	isResetOwner := resetOwner == user.Uuid

	if !users.IsEboard(user) && !users.IsActiveRTP(user) && !isCounterOwner && !isResetOwner {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := database.DeleteReset(ctx, counterID, req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error deleting reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	counterExists, err := database.UpdateCounterLastReset(ctx, counterID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error updating last reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, counterExists)
}

func Routes(r *gin.RouterGroup) {
	rGroup := r.Group("/reset")

	rGroup.GET("/:id")
	rGroup.POST("", CreateReset)
	rGroup.DELETE("", DeleteReset)

	rGroup.Handle("QUERY", "", GetResets)
}
