package modules

import (
	"testing"
	"time"
)

func TestPluginCore(t *testing.T){
	p:=PluginCore{HalPlugins: make(map[string]Plugin)}
	d:=Dummy{}
	p.Load("dummy", &d)
	p.PreInit()
	p.CallPlugin("dummy::Test","Wawa",17,time.Now())
}
