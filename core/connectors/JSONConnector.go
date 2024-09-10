package connectors

import (
	"HalogenGhostCore/core"
	"encoding/base64"
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
		cms := make([]core.CComment, 0)
		for _, comm := range comments {
			if r, err := base64.StdEncoding.DecodeString(comm.Comment); err == nil {
				comm.Comment = string(r)
			}
			cms = append(cms, comm)
		}

		c.output["comments"] = cms
		c.output["count"] = count
		c.output["page"] = page
	}
	c.Success("Comments retrieved")
}

func (c *JSONConnector) Comment_LevelGet(comments []core.CComment, count int, page int) {
	c.Comment_AccountGet(comments, count, page)
}

func (c *JSONConnector) Comment_HistoryGet(comments []core.CComment, acc core.CAccount, role core.Role, count int, page int) {
	c.Comment_AccountGet(comments, count, page)
	c.output["user"] = struct {
		ModBadge       int    `json:"mod_badge"`
		CommentColor   string `json:"comment_color"`
		Uname          string `json:"uname"`
		IconId         int    `json:"icon_id"`
		IconType       int    `json:"icon_type"`
		ColorPrimary   int    `json:"color_primary"`
		ColorSecondary int    `json:"color_secondary"`
		Special        int    `json:"special"`
	}{
		role.ModLevel,
		role.CommentColor,
		acc.Uname,
		acc.GetShownIcon(),
		acc.IconType,
		acc.ColorPrimary,
		acc.ColorSecondary,
		acc.Special,
	}
}

func (c *JSONConnector) Communication_FriendGetRequests(reqs []map[string]string, count int, page int) {
	c.output["requests"] = reqs
	c.output["count"] = count
	c.output["page"] = page
	c.Success("Friend requests retrieved")
}

func (c *JSONConnector) Communication_MessageGet(msg core.CMessage, uid int) {
	if content, err := base64.StdEncoding.DecodeString(msg.Message); err == nil {
		msg.Message = string(content)
	}
	uidx := msg.UidDest
	if uid == msg.UidDest {
		uidx = msg.UidSrc
	}
	xacc := core.CAccount{DB: msg.DB, Uid: uidx}
	xacc.LoadAuth(core.CAUTH_UID)
	c.output["content"] = struct {
		core.CMessage
		Uname string `json:"uname"`
	}{
		msg,
		xacc.Uname,
	}
	c.Success("Message retrieved")
}
