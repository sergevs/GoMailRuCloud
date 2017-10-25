package mailrucloud

import (
	"fmt"
	"io"
	//	"io/ioutil"
	"net/http"
	"os"
)

// Get downloads file parts or the full file from the cloud.
// src is the full file path in the cloud
// dst is the local path to store
// ch is progress channel
func (c *MailRuCloud) Get(src, dst string, ch chan<- int) (err error) {
	mpart := false
	_, err = c.Stat(src)
	if err != nil {
		_, err = c.Stat(fmt.Sprintf("%s.Multifile-Part%02d", src, 0))
		mpart = true
	}
	if err != nil {
		Logger.Println(err)
		return
	}
	f, err := os.Create(dst)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer f.Close()
	if mpart {
		for i := 0; ; i++ {
			if _, err = c.Stat(fmt.Sprintf("%s.Multifile-Part%02d", src, i)); err == nil {
				err = c.GetFilePart(f, fmt.Sprintf("%s.Multifile-Part%02d", src, i), ch)
				if err != nil {
					Logger.Println(err)
					return
				}
			} else {
				if i > 0 {
					err = nil
				}
				break
			}
		}
	} else {
		err = c.GetFilePart(f, src, ch)
	}
	return
}

// GetFilePart downloads file parts or the full file from the cloud.
// f is file descriptor
// src is the full file path in the cloud
// ch is progress channel
func (c *MailRuCloud) GetFilePart(f *os.File, src string, ch chan<- int) (err error) {
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
		err = fmt.Errorf("Get failed. Status: %d", r.StatusCode)
		Logger.Println(err)
		return
	}
	_, err = io.Copy(f, NewIoProgress(r.Body, int(r.ContentLength), ch))
	if err != nil {
		Logger.Println(err)
		return
	}
	return
}
