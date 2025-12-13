package client

import (
	"encoding/json"
	"mangahub/internal/tcp/models"
	"net"
)

type ProgressSync interface {
	Send(models.ProgressUpdate)
}

type ProgressSyncClient struct {
	conn net.Conn
	addr string
}

func NewProgressSyncClient(addr string) *ProgressSyncClient {
	return &ProgressSyncClient{addr: addr}
}

func (c *ProgressSyncClient) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *ProgressSyncClient) Send(update models.ProgressUpdate) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	data, _ := json.Marshal(update)
	data = append(data, '\n')

	_, err := c.conn.Write(data)
	return err
}
