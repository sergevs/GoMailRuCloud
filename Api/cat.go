package mailrucloud

import (
	"fmt"
	"io"
	//	"io/ioutil"
	"net/http"
	"os"
)

// Cat file at mail.ru cloud.
// src is the full file path.
func (c *MailRuCloud) Cat(src string) (err error) {
	if err = c.GetShardInfo(); err != nil {
		return
	}
	r, err := c.Client.Get(c.Shard.Get + src)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("Cat failed. Status: %d", r.StatusCode)
		Logger.Println(err)
		return
	} else {
	  _, err = io.Copy(os.Stdout, r.Body)
		if err != nil {
			Logger.Println(err)
			return
		}
	}
	return
}
