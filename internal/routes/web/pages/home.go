package pages

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
)

func GetHomePage(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   users.IsEboard(user),
		"RTP":      users.IsActiveRTP(user),
	})
}
