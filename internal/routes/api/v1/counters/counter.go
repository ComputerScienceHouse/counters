package counters

import (
	"context"
	"net/http"
	"time"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

func GetCounters(c *gin.Context) {
	_, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	defer cancel()

	var req models.GetCountersInput

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error binding req")
		c.Status(http.StatusBadRequest)
		return
	}

	counters, err := database.GetCounters(ctx, &req.Filter)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("error getting counters!")
		c.Status(http.StatusBadRequest)
		return
	}

	var hasMore bool = false
	var timestamp int64 = 0
	if len(counters) > 1 {
		hasMore = (len(counters) > (database.MAX_COUNTERS_IN_QUERY)-1)
		timestamp = counters[len(counters)-1].ResetTimestamp
	}

	c.JSON(http.StatusOK, models.GetCountersOutput{
		Counters:  counters,
		HasMore:   hasMore,
		Timestamp: &timestamp,
	})
}

func DeleteCounter(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting CSH auth last reset")
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	defer cancel()

	var req models.DeleteCounterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error casting request")
		c.Status(http.StatusBadRequest)
		return
	}

	counterOwner, err := database.GetCounterOwner(ctx, req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting counter owner")
		c.Status(http.StatusBadRequest)
		return
	}

	sameUser := counterOwner == user.Uuid
	if !users.IsEboard(user) && !users.IsActiveRTP(user) && !sameUser {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting user permissions")
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := database.DeleteCounter(ctx, req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error deleting counter")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func CreateCounter(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	defer cancel()

	var req models.CreateCounterWithResetInput

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		body, _ := c.Get(gin.BodyBytesKey)
		raw, _ := body.([]byte)

		logging.Logger.Warnf("bad request for creating counter! %s %s",
			err,
			string(raw),
		)
		c.Status(http.StatusBadRequest)
		return
	}

	counterID, err := database.CreateCounterWithReset(ctx, user.Uuid, user.Username, &req.Counter, &req.Reset)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counters", "error": err}).Warning("Error creating new database")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, counterID)
}

func Routes(r *gin.RouterGroup) {
	cGroup := r.Group("/counter")

	cGroup.GET("/:id")
	cGroup.POST("", CreateCounter)
	cGroup.DELETE("", DeleteCounter)

	cGroup.Handle("QUERY", "", GetCounters)
}
