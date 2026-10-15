package pages

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
)

func GetResetPage(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	idToUse := c.Param("id")

	c.HTML(http.StatusOK, "reset.tmpl", gin.H{
		"Id":       idToUse,
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   users.IsEboard(user),
		"RTP":      users.IsActiveRTP(user),
	})
}
