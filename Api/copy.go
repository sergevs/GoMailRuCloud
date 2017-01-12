package mailrucloud

import (
	"io/ioutil"
	"net/url"
	"path"
)

// Copy is convenient method to move files at the mail.ru cloud.
// src and dst should be the full source and destination file paths.
func (c *MailRuCloud) Copy(src, dst string) (err error) {
	tdir, err := ioutil.TempDir("/tmp", "goapi")
	if err != nil {
		return
	}
	if err = c.Mkdir(tdir); err != nil {
		return
	}
	if err = c.CopyA(src, tdir); err != nil {
		return
	}
	if err = c.Rename(tdir+"/"+path.Base(src), path.Base(dst)); err != nil {
		return
	}
	if err = c.CopyA(tdir+"/"+path.Base(dst), path.Base(dst)); err != nil {
		return
	}
	if err = c.Remove(tdir); err != nil {
		return
	}
	return
}

// CopyA method is the direct call to api url.
// It does not support rename. src is full source file path, targetDir is the directory to copy file to.
func (c *MailRuCloud) CopyA(src, targetDir string) (err error) {
	return c.postReq(c.url("file/move"),
		url.Values{
			"api":      {"2"},
			"conflict": {"strict"},
			"home":     {src},
			"folder":   {targetDir},
			"token":    {c.AuthToken},
		},
		"Copy")
}
