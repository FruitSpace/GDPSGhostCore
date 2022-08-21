# PLUGINS

### LevelPacks
- [ ] `onGauntletNew` - invoked when new gauntlet is created
- [ ] `onMapPackNew` - invoked when new map pack is created
### Communication
- To Be Done

### Essential

```go
package modules

// Plugin describes Plugin type; PluginCore is also Plugin
type Plugin interface {
	PreInit(*PluginCore, ...interface{})
	Unload(...interface{})
}

// PreInit Invoked to load anything
func PreInit(pch *PluginCore, args ...interface{})

// Unload Unloads everything
func Unload(args ...interface{})
```
### Player

```go
package modules

// OnPlayerNew Invoked when player is registered, but not yet activated account
func OnPlayerNew(uid int, uname string, email string)

// OnPlayerActivate Invoked when player first activated account
func OnPlayerActivate(uid int, uname string)

// OnPlayerLogin invoked when player commits login (regular, not gjp)
func OnPlayerLogin(uid int, uname string)

// OnPlayerBackup invoked when player uploads their backup
func OnPlayerBackup(uid int, decryptedBackup string)

//onPlayerSync is forbidden

// OnPlayerScoreUpdate invoked when player updates their score
func OnPlayerScoreUpdate(uid int, uname string, stats map[string]int)

var stats = map[string]int{
	"stars": 0,
	"diamonds": 0,
	"coins": 0,
	"ucoins": 0,
	"demons": 0,
	"cpoints": 0,
	"orbs": 0,
	"moons": 0,
	"special": 0,
	"lvlsCompleted": 0,
}
```
### Level

```go
package modules

import "strconv"

// OnLevelUpload invoked when level was uploaded
func OnLevelUpload(id int, name string, builder string, desc string)

// OnLevelUpdate invoked when level was updated
func OnLevelUpdate(id int, name string, builder string, desc string)

// OnLevelDelete invoked when level was deleted
func OnLevelDelete(id int, name string, builder string)

// OnLevelRate invoked when level was rated/rerated
func OnLevelRate(id int, name string, builder string, stars int, likes int, downloads int, length int, demonDiff int, isEpic bool, isFeatured bool, ratedBy map[string]string)

// OnLevelReport invoked when level was reported
func OnLevelReport(id int, name string, builder string, player string)

// OnLevelScore invoked when player published their score in level scoreboard
func OnLevelScore(id int, name string, player string, percent int, coins int)

var ratedBy = map[string]string{
	"uid":   strconv.Itoa(1),
	"uname": "Username",
}
```

### Communtication

```go
package modules

// Not implemented yet
```