package modules

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DiscordPacker struct {
	Chan *amqp.Channel
	Passive bool
	pch *PluginCore
}

func (mod *DiscordPacker) PreInit(pch *PluginCore, data ...interface{}){
	channel:=pch.CallPlugin("RabbitMQ::ConnChan")
	if channel[1].Interface()!=nil {
		mod.Passive=true
		return
	}
	rchan:=channel[0].Interface().(*amqp.Channel)
	mod.Chan=rchan
	rchan.QueueDeclare("gdps_bot",true,false,false,false,nil)
}

func (mod *DiscordPacker) GenPayload(t string, data map[string]string) string {
	b,_:=json.Marshal(struct {
		event string
		data map[string]string
	}{t,data})
	return string(b)
}

func (mod *DiscordPacker) OnPlayerActivate(uid int, uname string) {
	mod.pch.CallPlugin("RabbitMQ::PublishText",)
}

func (mod *DiscordPacker) Unload(...interface{}){}