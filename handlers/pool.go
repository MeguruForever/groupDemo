package handlers

import (
	"GroupDemo/models"
	"GroupDemo/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AddUser(c *gin.Context) {
	groupID := c.Param("id")
	userId := c.Query("user_id")

	pool := utils.Pool
	pool.Mu.Lock()
	defer pool.Mu.Unlock()

	if _, exists := pool.Pools[groupID]; !exists {
		pool.Pools[groupID] = &models.SingleGroupPool{
			Connections: make(map[string]*websocket.Conn),
			Broadcast:   make(chan models.Message),
		}
		//建立广播连接
		go pool.Pools[groupID].BroadcastMessages()
		fmt.Println("create group pool", groupID)
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		fmt.Println("websocket accept error:", err, userId)
		return
	}
	fmt.Println("add user", userId)
	pool.Pools[groupID].Connections[userId] = conn
	go handleMessages(groupID, userId, conn)
}
