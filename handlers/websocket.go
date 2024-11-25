package handlers

import (
	"GroupDemo/models"
	"GroupDemo/utils"
	"fmt"
	"github.com/gorilla/websocket"
)

func handleMessages(groupID, userID string, conn *websocket.Conn) {
	defer func() {
		utils.Pool.Mu.Lock()
		defer utils.Pool.Mu.Unlock()
		err := conn.Close()
		if err != nil {
			fmt.Println("websocket close error:", err)
			return
		}
		delete(utils.Pool.Pools[groupID].Connections, userID)
	}()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		fmt.Println("backend receive message", msg)
		if err != nil {
			break
		}
		msg.GroupID = groupID
		msg.UserID = userID
		msg.Time = msg.Time.Local()
		utils.Pool.Pools[groupID].Broadcast <- msg
		utils.Pool.Pools[groupID].SaveMessage(msg)
	}
}
