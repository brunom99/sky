package client

import (
	"encoding/json"
	"fmt"
	"golife/config"
	"golife/fly"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id                   string
	seed                 int64
	config               config.Config
	ws                   *websocket.Conn
	mutexWsClient        sync.Mutex
	aircrafts            []*fly.Aircraft
	mutexMapAircrafts    sync.Mutex
	randSeed             *rand.Rand
	disconnected         bool
	fnUpdateLastActivity func()
	totalAircrafts       int
}

type InfoClient struct {
	Seed           string `json:"seed"`
	TotalAircrafts int    `json:"total_aircrafts"`
}

func (c *Client) start() error {
	// procedural random
	randSource := rand.NewSource(c.seed)
	c.randSeed = rand.New(randSource)
	// send message to the client
	if err := c.sendMessageByWs(nil); err != nil {
		return err
	}
	// mutex lock
	c.mutexMapAircrafts.Lock()
	defer c.mutexMapAircrafts.Unlock()
	// create n aircrafts
	for i := 1; i <= c.config.Aircraft.Count; i += 1 {
		c.aircrafts = append(c.aircrafts, c.createAircraft())
	}
	return nil
}

func (c *Client) createAircraft() *fly.Aircraft {
	// channel in/out to communicate with aircraft
	chOut := make(chan fly.MessageFromAircraft)
	chIn := make(chan fly.MessageToAircraft)
	// create aircraft
	aircraft := fly.CreateAircraft(c.config, c.randSeed.Int63(), chOut, chIn)
	c.totalAircrafts++
	// for each aircraft, wait for message to the client
	go c.waitingAircraftMsg(chOut)
	// wake up aircraft
	go aircraft.WakeUp()
	return aircraft
}

func (c *Client) waitingAircraftMsg(chanFromAircraft chan fly.MessageFromAircraft) {
	for !c.disconnected { // avoid goroutine leaks
		// waiting for msg on aircraft chan
		msgFromAircraft := <-chanFromAircraft
		// total aircraft --
		if msgFromAircraft.Aircraft.IsFinish {
			c.totalAircrafts--
		}
		// transmit aircraft msg to the browser
		go func() {
			_ = c.sendMessageByWs(msgFromAircraft.Aircraft)
		}()
		// update last activity
		c.fnUpdateLastActivity()
	}
}

func (c *Client) terminateAircrafts() {
	// mutex lock
	c.mutexMapAircrafts.Lock()
	defer c.mutexMapAircrafts.Unlock()
	// send message to all aircrafts
	for _, aircraft := range c.aircrafts {
		aircraft.Terminate()
	}
	c.aircrafts = nil
	// client is disconnected
	c.disconnected = true
}

func (c *Client) sendMessageByWs(aircraft *fly.Aircraft) error {
	// lock websocket mutex
	c.mutexWsClient.Lock()
	defer c.mutexWsClient.Unlock()
	// message to the client
	msg := struct {
		Aircraft *fly.Aircraft `json:"aircraft"`
		Info     InfoClient    `json:"info"`
	}{
		aircraft,
		InfoClient{
			Seed:           fmt.Sprintf("%d", c.seed),
			TotalAircrafts: c.totalAircrafts,
		},
	}
	// send message
	if bytes, err := json.Marshal(msg); err == nil {
		return c.ws.WriteMessage(1, bytes)
	}
	return nil
}
