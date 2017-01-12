package mailrucloud

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
}

// NewCloud authenticates with mail.ru and returns a new object associated with user account.
// domain parameter should be "mail.ru"
func NewCloud(user, password, domain string) (*MailRuCloud, error) {
	var c MailRuCloud
	cookieJar, _ := cookiejar.New(nil)
	c.Client = &http.Client{Jar: cookieJar}
	if err := c.postReq("https://auth.mail.ru/cgi-bin/auth?lang=ru_RU&from=authpopup",
		url.Values{
			"page":     {"https://cloud.mail.ru/?from=promo"},
			"Domain":   {domain},
			"Login":    {user},
			"Password": {password},
		},
		"Auth"); err != nil {
		return nil, err
	}
	b, err := c.getReq(c.url("tokens/csrf"), nil, "Auth")
	if err != nil {
		return nil, err
	}
	var t AuthToken
	if err = json.Unmarshal(b, &t); err != nil {
		Logger.Println(err)
		return nil, err
	}
	c.AuthToken = t.Body.Token
	return &c, nil
}
