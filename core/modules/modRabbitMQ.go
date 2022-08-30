package modules

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	ID string
}

func (mod *RabbitMQ) PreInit(pch *PluginCore, data ...interface{}){}

func (mod *RabbitMQ) Unload(...interface{}){}

func (mod *RabbitMQ) ConnChan(host string, user string, password string) (*amqp.Channel,error) {
	var err error
	mod.Conn,err=amqp.Dial("amqp://"+user+":"+password+"@"+host+"/")
	if err!=nil{
		a:=amqp.Channel{}
		return  &a,err
	}
	return mod.Conn.Channel()
}

func (mod *RabbitMQ) Close(channel *amqp.Channel) {
	channel.Close()
	mod.Conn.Close()
}

func (mod *RabbitMQ) PublishText(channel *amqp.Channel, text string) error {
	return channel.Publish("","gdps_bot",false,false,amqp.Publishing{
		ContentType: "text/plain",
		Body: []byte(text),
	})
}