package handlers

import (
	"GroupDemo/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"GroupDemo/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroup(t *testing.T) {
	// 初始化Gin引擎
	r := gin.Default()
	r.POST("/groups", CreateGroup)

	// 创建一个模拟的HTTP请求
	group := models.Group{ID: "1", Name: "Test Group"}
	jsonValue, _ := json.Marshal(group)
	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 调用处理函数
	r.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusCreated, w.Code)

	// 检查响应体
	var responseGroup models.Group
	err := json.Unmarshal(w.Body.Bytes(), &responseGroup)
	assert.NoError(t, err)
	assert.Equal(t, group.ID, responseGroup.ID)
	assert.Equal(t, group.Name, responseGroup.Name)
}

func TestJoinGroup(t *testing.T) {
	// 初始化Gin引擎
	r := gin.Default()
	r.POST("/groups/:id/join", JoinGroup)

	// 创建一个模拟的HTTP请求
	groupID := "1"
	userID := "114514"
	req, _ := http.NewRequest("POST", "/groups/"+groupID+"/join?user_id="+userID, nil)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 调用处理函数
	r.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 检查响应体
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "joined", response["status"])

	// 检查用户是否成功加入群组
	utils.Mu.Lock()
	members := utils.GroupMembers[groupID]
	utils.Mu.Unlock()
	assert.Contains(t, members, userID)
}
