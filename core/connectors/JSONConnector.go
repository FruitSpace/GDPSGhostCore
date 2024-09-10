package connectors

import (
	"HalogenGhostCore/core"
	"encoding/json"
	"strconv"
)

type JSONConnector struct {
	output map[string]interface{}
}

func (c *JSONConnector) Output() string {
	d, _ := json.Marshal(c.output)
	return string(d)
}

func (c *JSONConnector) Error(code string, reason string) {
	c.output["status"] = "error"
	c.output["message"] = reason
	c.output["code"] = code
}

func (c *JSONConnector) Success(message string) {
	c.output["status"] = "success"
	c.output["message"] = message
}

func (c *JSONConnector) Account_Sync(savedata string) {
	c.output["savedata"] = savedata
	c.Success("Savedata present")
}

func (c *JSONConnector) Account_Login(uid int) {
	c.output["uid"] = strconv.Itoa(uid)
	c.Success("Logged in")
}

func (c *JSONConnector) Comment_AccountGet(comments []core.CComment, count int, page int) {
	if len(comments) == 0 {
		c.output["comments"] = []string{}
		c.output["count"] = 0
		c.output["page"] = page
	} else {
		c.output["comments"] = comments
		c.output["count"] = count
		c.output["page"] = page
	}
	c.Success("Comments retrieved")
}
