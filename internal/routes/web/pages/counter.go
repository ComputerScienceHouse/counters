package pages

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCounterPage(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	idToUse := c.Param("id")

	counter, err := database.GetCounterFromId(c.Request.Context(), idToUse)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter page", "error": err}).Warning("Error geting counter")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse())
		return
	}

	if counter == nil {
		c.Redirect(http.StatusNotFound, "/")
		return
	}

	c.HTML(http.StatusOK, "counter.tmpl", gin.H{
		"CounterID":           counter.CounterID,
		"CounterOwnerUUID":    counter.CounterOwner,
		"CounterTitle":        counter.CounterTitle,
		"CounterDescription":  counter.CounterDescription,
		"MostRecentTimestamp": counter.ResetOccuredAt,

		"Username": user.Username,
		"FullName": user.FullName,
		"UserID":   user.Uuid,

		"EBoard": users.IsEboard(user),
		"RTP":    users.IsActiveRTP(user),
	})
}
