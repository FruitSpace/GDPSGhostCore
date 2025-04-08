package modules

import (
	"fmt"
	"time"
)

type Dummy struct {
	pch *PluginCore
}

func (mod *Dummy) PreInit(pch *PluginCore, args ...interface{}){
	mod.pch=pch
	mod.pch.CallPlugin("dummy::Test","mako",190,time.Now())
}

func (mod *Dummy) Unload(...interface{}){}

func (mod *Dummy) Test(s string, i int, t time.Time) {
	fmt.Println("Got",s,"and",i,": ",t.Unix())
}