package mailrucloud

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"
)

const okfmt = "\t%-15s \u2713"
const failmrk = "\u2717"
const testdir = "/goapi_test/"
const testfile = "the_file.xyz"
const filesize = 1024 * 1024

func TestApi(t *testing.T) {
	c, err := NewCloud(os.ExpandEnv("$MAILRU_USER"), os.ExpandEnv("$MAILRU_PASSWORD"), os.ExpandEnv("$MAILRU_DOMAIN"))
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "NewCloud")

	err = c.Mkdir(testdir)
	if err != nil {
		t.Error(err, failmrk)
	}
	t.Logf(okfmt, "Mkdir")

	_, err = c.List(testdir)
	if err != nil {
		t.Error(err, failmrk)
	}
	t.Logf(okfmt, "List")

	f, err := os.Create(testfile)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Create file")

	fs, err := io.CopyN(f, rand.New(rand.NewSource(time.Now().UnixNano())), filesize)
	if err != nil || fs != filesize {
		t.Fatal(err, failmrk)
	}
	f.Close()
	t.Logf(okfmt, "Random fill")

	f, err = os.Open(testfile)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Open file")

	md5sum := md5.New()
	fs, err = io.Copy(md5sum, f)
	if err != nil || fs != filesize {
		t.Fatal(err, failmrk)
	}
	f.Close()
	md5orig := md5sum.Sum(nil)
	t.Logf(okfmt+" (%x)", "Calc md5sum ", md5orig)

	err = c.Upload(testfile, testdir+testfile, nil)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Uploadfile")

	err = os.Remove(testfile)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Remove file")

	err = c.Get(testdir+testfile, testfile, nil)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Get")

	f, err = os.Open(testfile)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Open file")

	md5sum = md5.New()
	fs, err = io.Copy(md5sum, f)
	if err != nil || fs != filesize {
		t.Fatal(err, failmrk)
	}
	f.Close()
	md5get := md5sum.Sum(nil)
	t.Logf(okfmt+" (%x)", "Calc md5sum ", md5get)

	if fmt.Sprintf("%x", md5orig) != fmt.Sprintf("%x", md5get) {
		t.Errorf("MD5sum differs ! %s", failmrk)
	}
	t.Logf(okfmt, "MD5 sum match")

	err = c.Remove(testdir)
	if err != nil {
		t.Error(err, failmrk)
	}
	t.Logf(okfmt, "Remove")

	err = os.Remove(testfile)
	if err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Remove file")
}
