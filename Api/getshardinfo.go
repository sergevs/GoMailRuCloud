package mailrucloud

import (
	"encoding/json"
	"net/url"
)

// GetShardInfo retrives and unmarshal information about get and upload urls.
// There are many others urls but at now only that 2 are in use.
func (c *MailRuCloud) GetShardInfo() (err error) {
	b, err := c.getReq(c.url("dispatcher?"),
		url.Values{"token": {c.AuthToken}},
		"GetShardInfo")
	if err != nil {
		return
	}
	var s ShardInfo
	err = json.Unmarshal(b, &s)
	if err != nil {
		Logger.Println(err)
		return
	}
	c.Shard.Get = s.Body.Get[0].Url
	c.Shard.Upload = s.Body.Upload[0].Url
	return
}
