package usercontrollers

import (
	"database-api/database/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteUserController(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	deleterUser, err := user.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deleterUser)
}
