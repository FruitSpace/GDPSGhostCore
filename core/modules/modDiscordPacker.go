package modules

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
)

type DiscordPacker struct {
	Conn *amqp.Connection
	ID string
	Passive bool
}

func (mod *DiscordPacker) PreInit(pch *PluginCore, data ...interface{}){
	channel:=pch.CallPlugin("RabbitMQ::ConnChan")
	if channel[1].Interface()!=nil {
		mod.Passive=true
		return
	}
	rchan:=channel[0].Convert(reflect.TypeOf(amqp.Channel{}))
	rchan.QueueDeclare("bot_"+mod.ID,true,false,false,false,nil)
}

func (mod *DiscordPacker) Unload(...interface{}){}