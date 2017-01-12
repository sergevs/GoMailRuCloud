package main

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/sergevs/GoMailRuCloud/Api"
	"log"
	"os"
	"path"
	"path/filepath"
)

func progress(c <-chan int) {
	uiprogress.Start()
	b := uiprogress.AddBar(<-c)
	for i := range c {
		b.Set(i)
	}
	uiprogress.Stop()
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s is Trivial Mail.ru cloud Client\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Usage: %s -<COMMAND> [FILE/DIR] [FILE/DIR]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Provides basic operates with files at mail.ru cloud storage\n\n")
	fmt.Fprintf(os.Stderr, "COMMAND := < cp | cat | get | ls | mkdir | mv | put | rm | stat >\n\n")
	fmt.Fprintf(os.Stderr, "Example: tmrc -ls\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	//  Uncomment the next line to disable very useful default logger
	//  mailrucloud.Logger = log.New(ioutil.Discard, "", log.Lshortfile)
	c, err := mailrucloud.NewCloud(os.ExpandEnv("$MAILRU_USER"), os.ExpandEnv("$MAILRU_PASSWORD"), "mail.ru")
	if err != nil {
		os.Exit(1)
	}
	switch cmd := os.Args[1]; cmd {
	case "-ls":
		var dir string
		if len(os.Args) > 2 {
			dir = os.Args[2]
		} else {
			dir = "/"
		}
		err = c.PrintFileList(dir)
	case "-rm":
		if len(os.Args) > 2 {
			err = c.Remove(os.Args[2])
		} else {
			log.Fatal("File or dir is not specified")
		}
	case "-mkdir":
		if len(os.Args) > 2 {
			err = c.Mkdir(os.Args[2])
		} else {
			log.Fatal("Dir is not specified")
		}
	case "-stat":
		if len(os.Args) > 2 {
			err = c.PrintFileStat(os.Args[2])
		} else {
			log.Fatal("Dir is not specified")
		}
	case "-cat":
		if len(os.Args) != 3 {
			log.Fatal("not enougth arguments")
		}
 		err = c.Cat(os.Args[2])
	case "-get":
		var dst string
		if len(os.Args) > 3 {
			dst = os.Args[3]
		} else if len(os.Args) < 3 {
			log.Fatal("not enougth arguments")
		} else {
			dst = path.Base(os.Args[2])
		}
		pc := make(chan int)
		go progress(pc)
		err = c.Get(os.Args[2], dst, pc)
	case "-put":
		var dst string
		if len(os.Args) > 3 {
			dst = os.Args[3]
		} else if len(os.Args) < 3 {
			log.Fatal("not enougth arguments")
		} else {
			dst = "/" + filepath.Base(os.Args[2])
		}
		pc := make(chan int)
		go progress(pc)
		err = c.Upload(os.Args[2], dst, pc)
	case "-mv":
		if len(os.Args) > 3 {
			err = c.Move(os.Args[2], os.Args[3])
		} else {
			log.Fatal("not enougth arguments")
		}
	case "-cp":
		if len(os.Args) > 3 {
			err = c.Copy(os.Args[2], os.Args[3])
		} else {
			log.Fatal("not enougth arguments")
		}
	default:
		usage()
		log.Fatal("Wrong command")
	}
	if err != nil {
		os.Exit(1)
	}
}
