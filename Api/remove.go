package mailrucloud

import (
	"net/url"
)

// Remove a file at the mail.ru cloud.
func (c *MailRuCloud) Remove(path string) (err error) {
	return c.postReq(c.url("file/remove"),
		url.Values{
			"api":   {"2"},
			"home":  {path},
			"token": {c.AuthToken},
		},
		"Remove")
}
