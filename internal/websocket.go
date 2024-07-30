package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
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
	crawl ICrawl
	conn  *websocket.Conn
	send  chan *Message
}

var (
	clients = make(map[ICrawl]*Client)
	mu      sync.Mutex
)

func newClient(crawl ICrawl, ctx *gin.Context) (*Client, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return &Client{crawl: crawl, conn: conn, send: make(chan *Message)}, nil
}

func (c *Client) Read() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("error: %v", err)
	}

	c.conn.SetPongHandler(func(string) error {
		err = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("error: %v", err)

			return err
		}
		return nil
	})

	for {
		dataType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch dataType {
		case websocket.TextMessage:
			fmt.Println("Received Text Message:", string(data))

			//str := string(data)
			message := Message{}
			//if str == "init" {
			message.Cmd = CrawlCmd_Login
			//} else if str == "login" {
			//	message.Cmd = CrawlCmd_Login
			//}

			c.crawl.GetChan() <- &ClientMessage{message: &message, client: c}

		//case websocket.BinaryMessage:
		//	fmt.Println("Received Binary Message:", data)
		//
		//	message := &Message{}
		//
		//	if err := proto.Unmarshal(data, message); err != nil {
		//		log.Printf("unmarshal data error:%v\n", err)
		//		continue
		//	}
		//
		//	c.crawl.GetChan() <- message
		default:
			fmt.Println("Received Unknown Message Type")
		}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}

			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message.Data)
			if err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func HandleWebsocket(accountId int, ctx *gin.Context) error {
	c, err := GetCrawl(accountId)
	if err != nil {
		return err
	}

	mu.Lock()
	client, err := newClient(c, ctx)
	if err != nil {
		return err
	}

	if oldClient, ok := clients[c]; ok {
		err := oldClient.conn.Close()
		if err != nil {
			return err
		}
	}

	clients[c] = client
	mu.Unlock()

	go client.Read()
	go client.Write()

	return nil
}
