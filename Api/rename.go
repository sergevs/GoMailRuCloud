package mailrucloud

import (
	"net/url"
)

// Rename a file at the mail.ru cloud.
func (c *MailRuCloud) Rename(src, targetName string) (err error) {
	return c.postReq(c.url("file/rename"),
		url.Values{
			"api":      {"2"},
			"conflict": {"strict"},
			"home":     {src},
			"name":     {targetName},
			"token":    {c.AuthToken},
		},
		"Rename")
}
