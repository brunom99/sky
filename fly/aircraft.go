package fly

import (
	"golife/config"
	"golife/utils"
	"math/rand"
	"time"
)

type Aircraft struct {
	ID           string   `json:"id"`
	Position     Position `json:"pos"`
	IsFinish     bool     `json:"is_finish"`
	PositionPrev Position `json:"-"`
	randSeed     *rand.Rand
	chOut        chan MessageFromAircraft
	chIn         chan MessageToAircraft
	config       config.Config
	speed        int
	distance     float64
	targetPos    Position
}

func CreateAircraft(config config.Config, seed int64,
	chOut chan MessageFromAircraft, chIn chan MessageToAircraft) *Aircraft {
	randSeed := rand.New(rand.NewSource(seed))
	posStart := RandPosition(randSeed)
	aircraft := Aircraft{
		ID:           utils.Uuid(),
		Position:     posStart,
		PositionPrev: posStart,
		randSeed:     randSeed,
		chOut:        chOut,
		chIn:         chIn,
		config:       config,
		targetPos:    posStart,
	}
	aircraft.speed = utils.RandInt(config.Aircraft.MinResponse, config.Aircraft.MaxResponse, randSeed)
	aircraft.distance = utils.RandFloat(config.Aircraft.MinDistance, config.Aircraft.MaxDistance, randSeed)
	return &aircraft
}

func (a *Aircraft) WakeUp() {
	// check channel: message to aircraft
	go a.readMessage()
	// while aircraft is alive
	for !a.IsFinish { // avoid goroutine leaks
		// send message to the client
		a.sendMessage()
		// random waiting
		time.Sleep(time.Duration(a.speed) * time.Millisecond)
		// move aircraft
		a.move()
	}
	// aircraft is finish: close read channel
	close(a.chIn)
	// aircraft is finish: send a last message
	a.sendMessage()
}

func (a *Aircraft) Message(msg MessageToAircraft) {
	a.chIn <- msg
}

func (a *Aircraft) Terminate() {
	a.IsFinish = true
}

func (a *Aircraft) sendMessage() {
	a.chOut <- MessageFromAircraft{
		Aircraft: a,
	}
}

func (a *Aircraft) readMessage() {
	// while aircraft is alive
	for !a.IsFinish {
		// wait for channel message
		message := <-a.chIn
		_ = message
	}
}

func (a *Aircraft) move() {
	// define new targetPos ?
	if a.Position.IsSame(a.targetPos) {
		a.targetPos = RandPosition()
	}
	// save old position
	a.PositionPrev = a.Position
	// move
	if !a.Position.Move(a.targetPos, a.distance) {
		// target became current position
		a.targetPos = a.Position
	}
}
