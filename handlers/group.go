package handlers

import (
	"GroupDemo/models"
	"GroupDemo/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.Mu.Lock()
	utils.Groups[group.ID] = group.Name
	fmt.Println(group.Name)
	utils.Mu.Unlock()
	c.JSON(http.StatusCreated, group)
}

func JoinGroup(c *gin.Context) {
	groupID := c.Param("id")
	userID := c.Query("user_id")
	utils.Mu.Lock()
	utils.GroupMembers[groupID] = append(utils.GroupMembers[groupID], userID)
	utils.Mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"status": "joined"})
}
