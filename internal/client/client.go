package client

import (
	"encoding/json"
	"fmt"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/message"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	conn        *websocket.Conn
	account     *account.Account
	accountChan chan *message.AsyncMessage
	done        chan struct{}
}

var (
	clients = make(map[*account.Account]*Client)
	mu      sync.Mutex
)

func newClient(id int, ctx *gin.Context) (*Client, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:        conn,
		account:     nil,
		accountChan: make(chan *message.AsyncMessage),
		done:        make(chan struct{}),
	}

	return client, nil
}

func (c *Client) readPump() {
	defer func() {
		close(c.done)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		// 解析消息并处理
		msg := &message.Message{}
		if err := json.Unmarshal(data, msg); err != nil {
			log.Printf("error parsing data: %v", err)
			continue
		}

		if err := c.handleWebSocketMessage(msg); err != nil {
			log.Printf("error handling websocket data: %v", err)
			continue
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg := <-c.accountChan:
			if err := c.handleAccountMessage(msg); err != nil {
				log.Printf("handle account message error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *Client) handleWebSocketMessage(msg *message.Message) error {
	switch msg.Cmd {
	case message.CrawlCmd_Login:
		// 创建任务消息
		asyncMessage := message.NewAsyncMessage(msg, nil)

		// 发送任务并等待结果
		c.accountChan <- asyncMessage
		result := <-asyncMessage.Result

		// 发送结果回 WebSocket
		msg := &message.Message{
			Cmd:  message.CrawlCmd_Login,
			Data: []byte(fmt.Sprintf("%v", result)),
		}
		c.accountChan <- message.NewAsyncMessage(msg, nil)
	}
	return nil
}

func (c *Client) handleAccountMessage(msg *message.AsyncMessage) error {
	// 将消息写入 WebSocket
	return c.conn.WriteMessage(websocket.TextMessage, msg.Data)
}

func (c *Client) Close() {
	// 关闭 done channel 来通知所有 goroutine 退出
	close(c.done)

	// 关闭 WebSocket 连接
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.conn.Close()

	// 关闭 accountChan
	close(c.accountChan)

	// 从全局 clients map 中移除
	mu.Lock()
	delete(clients, c.account)
	mu.Unlock()
}

func HandleWebsocket(accountId int, ctx *gin.Context) error {
	//account, err := GetCrawl(accountId)
	//if err != nil {
	//	return err
	//}
	//
	//mu.Lock()
	//client, err := newClient(account, ctx)
	//if err != nil {
	//	mu.Unlock()
	//	return err
	//}
	//
	//if oldClient, ok := clients[account]; ok {
	//	oldClient.Close()
	//}
	//
	//clients[account] = client
	//mu.Unlock()
	//
	//go client.readPump()
	//go client.writePump()

	return nil
}