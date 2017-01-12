package mailrucloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *MailRuCloud) url(o string) string {
	return "https://cloud.mail.ru/api/v2/" + o
}

func (c *MailRuCloud) postReq(Url string, data url.Values, errPrefix string) (err error) {
	r, err := c.Client.PostForm(Url, data)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Println(err)
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s failed. Status: %d, Msg: %s", errPrefix, r.StatusCode, string(b))
		Logger.Println(err)
		return
	}
	return
}

func (c *MailRuCloud) getReq(Url string, data url.Values, errPrefix string) (b []byte, err error) {
	if data != nil {
		Url = Url + data.Encode()
	}
	r, err := c.Client.Get(Url)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Println(err)
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s failed. Status: %d, Msg: %s", errPrefix, r.StatusCode, string(b))
		Logger.Println(err)
		return
	}
	return
}
