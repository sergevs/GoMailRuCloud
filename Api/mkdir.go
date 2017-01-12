package mailrucloud

import (
	"net/url"
)

// Mkdir creates directory in the mail.ru cloud
func (c *MailRuCloud) Mkdir(path string) (err error) {
	return c.postReq(c.url("folder/add"),
		url.Values{
			"api":      {"2"},
			"conflict": {"strict"},
			"home":     {path},
			"token":    {c.AuthToken},
		},
		"Mkdir")
}
