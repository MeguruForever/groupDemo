package handlers

import (
	"GroupDemo/models"
	"GroupDemo/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestAddUser(t *testing.T) {
	// 初始化Gin引擎
	r := gin.Default()
	r.GET("/groups/:id/join", AddUser)

	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(r)
	defer server.Close()

	// 创建一个WebSocket连接
	groupID := "1"
	userID := "114514"
	wsURL := "ws" + server.URL[4:] + "/groups/" + groupID + "/join?user_id=" + userID
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// 检查用户是否成功加入群组
	utils.Pool.Mu.Lock()
	defer utils.Pool.Mu.Unlock()
	_, exists := utils.Pool.Pools[groupID].Connections[userID]
	assert.True(t, exists)
}

func TestHandleMessages(t *testing.T) {
	// Initialize Gin engine
	r := gin.Default()
	r.GET("/groups/:id/join", AddUser)

	// Create a mock HTTP server
	server := httptest.NewServer(r)
	defer server.Close()

	// Create WebSocket connections for two users
	groupID := "1"
	userID1 := "114514"
	userID2 := "1919810"
	wsURL1 := "ws" + server.URL[4:] + "/groups/" + groupID + "/join?user_id=" + userID1
	wsURL2 := "ws" + server.URL[4:] + "/groups/" + groupID + "/join?user_id=" + userID2
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	assert.NoError(t, err)
	defer ws1.Close()
	ws2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	assert.NoError(t, err)
	defer ws2.Close()

	// Send a message from user1
	message := models.Message{
		Content: "Hello, World! 我是114514",
	}
	err = ws1.WriteJSON(message)
	assert.NoError(t, err)

	// Read the message from user2's connection
	var receivedMessage models.Message
	err = ws2.ReadJSON(&receivedMessage)
	fmt.Println("received", receivedMessage)
	assert.NoError(t, err)
	assert.Equal(t, message.Content, receivedMessage.Content)

	// Verify the message is saved on the backend
	utils.Pool.Mu.Lock()
	defer utils.Pool.Mu.Unlock()
	assert.Contains(t, utils.Pool.Pools[groupID].Connections, userID1)
	assert.Contains(t, utils.Pool.Pools[groupID].Connections, userID2)
}
