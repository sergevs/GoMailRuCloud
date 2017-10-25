package mailrucloud

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// mail.ru limits file size to 2GB
// the current implementation is limited to max int ( 4 bytes ) size of bytes.Buffer
// which used to send multipart form , 1024 bytes reserved for form data/fields
// const MaxFileSize = 3

const MaxFileSize = 2*1024*1024*1024 - 1024

// Upload is a convenient method to upload files to the mail.ru cloud.
// src is the local file path
// dst is the full destination file path
// ch  is a channel to report operation progress. can be nil.
func (c *MailRuCloud) Upload(src, dst string, ch chan<- int) (err error) {
	if err = c.GetShardInfo(); err != nil {
		return
	}
	f, err := os.Open(src)
	if err != nil {
		return
	}
	s, err := f.Stat()
	if err != nil {
		return
	}
	if s.Size() <= MaxFileSize {
		return c.UploadFilePart(f, dst, 0, s.Size(), ch)
	} else {
		for spos, part := int64(0), 0; spos < s.Size(); spos, part = spos+MaxFileSize, part+1 {
			var n int64
			if spos+MaxFileSize <= s.Size() {
				n = MaxFileSize
			} else {
				n = s.Size() % MaxFileSize
			}
			//      fmt.Printf("spos %d %d %d\n", spos, s.Size(),n )
			if err = c.UploadFilePart(f, fmt.Sprintf("%s.Multifile-Part%02d", dst, part), spos, n, ch); err != nil {
				return
			}
		}
		return
	}
}

// UploadFilePart is the method to overcome 2Gb file size limit.
// file is the uploaded file descriptor
// dst is the full destination file path. the large files will be splitted appending .partXX extention
// spos is starting position of file
// n is number of bytes to write
// ch  is a channel to report operation progress. can be nil.
func (c *MailRuCloud) UploadFilePart(file *os.File, dst string, spos, n int64, ch chan<- int) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.SetBoundary(randomBoundary())
	fw, err := w.CreateFormFile("file", "filename")
	if err != nil {
		Logger.Println(err)
		return
	}
	if _, err = io.Copy(fw, io.NewSectionReader(file, spos, n)); err != nil {
		Logger.Println(err)
		return
	}
	w.Close()
	Url := c.Shard.Upload
	req, err := http.NewRequest("POST", Url, NewIoProgress(&b, b.Len(), ch))
	if err != nil {
		Logger.Println(err)
		return
	}
	req.ContentLength = int64(b.Len())
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept-Encoding", "*.*")

	//	dump, _ := httputil.DumpRequestOut(req, true)
	//	fmt.Printf("%q", dump)
	r, err := c.Client.Do(req)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	// Check the response
	br, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Println(err)
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("Put file failed. Status: %d, Msg: %s", r.StatusCode, string(br))
		Logger.Println(err)
		return
	}
	hs := strings.SplitN(strings.TrimSpace(string(br)), ";", 2)
	err = c.addFile(dst, hs[0], hs[1])
	return
}

func (c *MailRuCloud) addFile(dst, hash, size string) (err error) {
	Url := c.url("file/add")
	data := url.Values{
		"token":    {c.AuthToken},
		"home":     {dst},
		"conflict": {"strict"},
		"hash":     {hash},
		"size":     {size},
	}
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
		err = fmt.Errorf("addFile failed. Status: %d, Msg: %s", r.StatusCode, string(b))
		Logger.Println(err)
		return
	}
	return
}

// the default function return too long boundary
// mail.ru does not accept it
func randomBoundary() string {
	var buf [15]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}
