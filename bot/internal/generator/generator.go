package generator

import (
	"fmt"
)

type MessageGenerator struct {
	counter    int
	userModulo int
}

func NewMessageGenerator(userCount int) *MessageGenerator {
	if userCount <= 0 {
		userCount = 3
	}
	return &MessageGenerator{
		counter:    1,
		userModulo: userCount,
	}
}

func (g *MessageGenerator) Next() (userID string, content string) {
	g.counter++
	userID = fmt.Sprintf("user-%d", g.counter%g.userModulo)
	content = fmt.Sprintf("This is message number %d", g.counter)
	return
}
