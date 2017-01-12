package mailrucloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// Stat is a convenient method to get information about a file at mail.ru cloud.
func (c *MailRuCloud) Stat(path string) (finfo *FileList, err error) {
	b, err := c.getReq(c.url("file?"),
		url.Values{"home": {path},
			"token": {c.AuthToken},
		},
		"Stat")
	if err != nil {
		return
	}
	var f FileList
	err = json.Unmarshal(b, &f)
	if err != nil {
		Logger.Println(err)
		return
	}
	finfo = &f
	return
}

// PrintFileStat is a convenient method to print information about a file at mail.ru cloud.
func (c *MailRuCloud) PrintFileStat(path string) (err error) {
	s, err := c.Stat(path)
	if err != nil {
		return
	}
	i := s.Body
	t := time.Unix(int64(i.Mtime), 0)
	var ft string
	if time.Now().Unix()-t.Unix() < 31536000 {
		ft = fmt.Sprintf("%.3s %2d %02d:%02d", t.Month(), t.Day(), t.Hour(), t.Minute())
	} else {
		ft = fmt.Sprintf("%.3s %2d %5d", t.Month(), t.Day(), t.Year())
	}
	fmt.Printf("%-6s %11d %-12s %-13s\n", i.Kind, i.Size, ft, i.Name)
	return
}
