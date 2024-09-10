// Package modules is an ultimate plugin core with plugin autoload
package modules

import (
	"HalogenGhostCore/core"
	"reflect"
	"strings"
)

type Plugin interface {
	PreInit(*PluginCore, ...interface{})
	Unload(...interface{})
}

type PluginCore struct {
	HalPlugins map[string]Plugin
}

func (pch *PluginCore) LoadPrepared(conf core.ConfigBlob) {
	//for mod, en := range conf.ServerConfig.EnableModules {
	//	pch.Load(mod, en)
	//}
}

func (pch *PluginCore) Load(name string, plugin Plugin) {
	pch.HalPlugins[name] = plugin
}

func (pch *PluginCore) CallPlugin(endpoint string, args ...interface{}) []reflect.Value {
	_endpoint := strings.Split(endpoint, "::") // PluginName::Method
	if plug, ok := pch.HalPlugins[_endpoint[0]]; ok {
		if _, ok := reflect.TypeOf(plug).MethodByName(_endpoint[1]); ok {
			//if plugin exists and has a method then convert all data to reflect.Value, call method and return its output
			inputs := make([]reflect.Value, 0, len(args)+1)
			for i := range args {
				inputs = append(inputs, reflect.ValueOf(args[i]))
			}
			return reflect.ValueOf(plug).MethodByName(_endpoint[1]).Call(inputs)
		}
	}
	return []reflect.Value{}
}

//===ESSENTIAL===

// PreInit Invoked to load anything
func (pch *PluginCore) PreInit(args ...interface{}) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::PreInit", pch, args)
	}
}

// Unload Unloads everything
func (pch *PluginCore) Unload(args ...interface{}) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::Unload", args)
	}
}

//===PLAYER===

// OnPlayerNew Invoked when player is registered, but not yet activated account
func (pch *PluginCore) OnPlayerNew(uid int, uname string, email string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnPlayerNew", uid, uname, email)
	}
}

// OnPlayerActivate Invoked when player first activated account
func (pch *PluginCore) OnPlayerActivate(uid int, uname string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnPlayerActivate", uid, uname)
	}
}

// OnPlayerLogin invoked when player commits login (regular, not gjp)
func (pch *PluginCore) OnPlayerLogin(uid int, uname string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnPlayerLogin", uid, uname)
	}
}

// OnPlayerBackup invoked when player uploads their backup
func (pch *PluginCore) OnPlayerBackup(uid int, decryptedBackup string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnPlayerBackup", uid, decryptedBackup)
	}
}

// OnPlayerScoreUpdate invoked when player updates their score
func (pch *PluginCore) OnPlayerScoreUpdate(uid int, uname string, stats map[string]int) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnPlayerScoreUpdate", uid, uname, stats)
	}
}

//===LEVEL===

// OnLevelUpload invoked when level was uploaded
func (pch *PluginCore) OnLevelUpload(id int, name string, builder string, desc string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnLevelUpload", id, name, builder, desc)
	}
}

// OnLevelUpdate invoked when level was updated
func (pch *PluginCore) OnLevelUpdate(id int, name string, builder string, desc string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnLevelUpdate", id, name, builder, desc)
	}
}

// OnLevelDelete invoked when level was deleted
func (pch *PluginCore) OnLevelDelete(id int, name string, builder string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnLevelDelete", id, name, builder)
	}
}

// OnLevelRate invoked when level was rated/rerated
func (pch *PluginCore) OnLevelRate(id int, name string, builder string, stars int, likes int, downloads int, length int, demonDiff int, isEpic bool, isFeatured bool, ratedBy map[string]string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnLevelRate", id, name, builder, stars, likes, downloads, length, demonDiff, isEpic, isFeatured, ratedBy)
	}
}

// OnLevelReport invoked when level was reported
func (pch *PluginCore) OnLevelReport(id int, name string, builder string, player string) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnLevelReport", id, name, builder, player)
	}
}

// OnLevelScore invoked when player published their score in level scoreboard
func (pch *PluginCore) OnLevelScore(id int, name string, player string, percent int, coins int) {
	for plug := range pch.HalPlugins {
		pch.CallPlugin(plug+"::OnLevelScore", id, name, player, percent, coins)
	}
}
