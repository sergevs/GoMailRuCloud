package mailrucloud

import (
	"io/ioutil"
	"net/url"
	"path"
)

// Move is convenient method to move files in the mail.ru cloud.
// src and dst should be the full source and destination file paths.
func (c *MailRuCloud) Move(src, dst string) (err error) {
	if path.Dir(src) == path.Dir(dst) {
		return c.Rename(src, path.Base(dst))
	}
	if path.Base(src) == path.Base(dst) {
		return c.MoveA(src, path.Dir(dst))
	}
	tdir, err := ioutil.TempDir("/tmp", "goapi")
	if err != nil {
		return
	}
	if err = c.Mkdir(tdir); err != nil {
		return
	}
	if err = c.MoveA(src, tdir); err != nil {
		return
	}
	if err = c.Rename(tdir+"/"+path.Base(src), path.Base(dst)); err != nil {
		return
	}
	if err = c.MoveA(tdir+"/"+path.Base(dst), path.Dir(dst)); err != nil {
		return
	}
	if err = c.Remove(tdir); err != nil {
		return
	}
	return
}

// MoveA method is the direct call to api url.
// It does not support rename. src is full source file path, targetDir is the directory to move file to.
func (c *MailRuCloud) MoveA(src, targetDir string) (err error) {
	return c.postReq(c.url("file/move"),
		url.Values{
			"api":      {"2"},
			"conflict": {"strict"},
			"home":     {src},
			"folder":   {targetDir},
			"token":    {c.AuthToken},
		},
		"Move")
}
