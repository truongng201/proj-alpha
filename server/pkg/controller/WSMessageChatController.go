package controller

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSMessageChatController struct {}


func (controller WSMessageChatController) Execute(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	fmt.Println("Client Message Chat connected:", ws.RemoteAddr())

	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("Message Chat Received : %s\n", msg)

		// Write
		err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			c.Logger().Error(err)
		}

	}
}