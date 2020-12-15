package hci

import (
	"bytes"
	"sync"
)

// Pool ...
type Pool struct {
	sync.Mutex

	sz   int
	cnt  int
	ch   chan *bytes.Buffer
	sent chan *bytes.Buffer
}

// NewPool ...
func NewPool(sz int, cnt int) *Pool {
	ch := make(chan *bytes.Buffer, cnt)
	sent := make(chan *bytes.Buffer, cnt)
	for len(ch) < cnt {
		ch <- bytes.NewBuffer(make([]byte, sz))
	}
	return &Pool{sz: sz, cnt: cnt, ch: ch, sent: sent}
}

// Client ...
type Client struct {
	p *Pool
}

// NewClient ...
func NewClient(p *Pool) *Client {
	return &Client{p: p}
}

// LockPool ...
func (c *Client) LockPool() {
	c.p.Lock()
}

// UnlockPool ...
func (c *Client) UnlockPool() {
	c.p.Unlock()
}

// Get returns a buffer from the shared buffer pool.
func (c *Client) Get() *bytes.Buffer {
	b := <-c.p.ch
	b.Reset()
	c.p.sent <- b
	return b
}

// Put puts the oldest sent buffer back to the shared pool.
func (c *Client) Put() {
	select {
	case b := <-c.p.sent:
		c.p.ch <- b
	default:
	}
}

// PutAll puts all the sent buffers back to the shared pool.
func (c *Client) PutAll() {
	for {
		select {
		case b := <-c.p.sent:
			c.p.ch <- b
		default:
			return
		}
	}
}
