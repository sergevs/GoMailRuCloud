package mailrucloud

import (
	"fmt"
	"io"
	//	"io/ioutil"
	"net/http"
	"os"
)

// Get downloads file from the cloud.
// src is the full file path in the cloud, dst is the local path to store.
func (c *MailRuCloud) Get(src, dst string, ch chan<- int) (err error) {
	if err = c.GetShardInfo(); err != nil {
		return
	}
	r, err := c.Client.Get(c.Shard.Get + src)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	var f *os.File
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("Get failed. Status: %d", r.StatusCode)
		Logger.Println(err)
		return
	} else {
		f, err = os.Create(dst)
		if err != nil {
			Logger.Println(err)
			return
		}
	}
	_, err = io.Copy(f, NewIoProgress(r.Body, int(r.ContentLength), ch))
	if err != nil {
		Logger.Println(err)
		return
	}
	return
}
