package mailrucloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// List the files at the mail.ru cloud.
func (c *MailRuCloud) List(path string) (list *FileList, err error) {
	b, err := c.getReq(c.url("folder?"),
		url.Values{"home": {path},
			"token": {c.AuthToken},
		},
		"List")
	if err != nil {
		return
	}
	var f FileList
	err = json.Unmarshal(b, &f)
	if err != nil {
		Logger.Println(err)
		return
	}
	list = &f
	return
}

// PrintFileList is a convenient method method to prin files list at the mail.ru cloud.
func (c *MailRuCloud) PrintFileList(path string) (err error) {
	l, err := c.List(path)
	if err != nil {
		return
	}
	for _, i := range l.Body.List {
		t := time.Unix(int64(i.Mtime), 0)
		var ft string
		if time.Now().Unix()-t.Unix() < 31536000 {
			ft = fmt.Sprintf("%.3s %2d %02d:%02d", t.Month(), t.Day(), t.Hour(), t.Minute())
		} else {
			ft = fmt.Sprintf("%.3s %2d %5d", t.Month(), t.Day(), t.Year())
		}
		fmt.Printf("%-6s %11d %-12s %-13s\n", i.Kind, i.Size, ft, i.Name)
	}
	return
}
