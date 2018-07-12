package models

import (
	"encoding/json"
	"time"
)

type ActionType string

const (
	FORWARD  = ActionType("forward")
	RECEIVED = ActionType("received")
)

type DeliveryModel struct {
	Date      time.Time
	Signature string
	Action    ActionType
	Data      interface{}
}

func (d *DeliveryModel) GetForwardModel() ForwardModel {
	b, _ := json.Marshal(d.Data)
	i := ForwardModel{}
	json.Unmarshal(b, &i)
	return i
}
