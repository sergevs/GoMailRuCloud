package mailrucloud

import (
	"net/http"
)

// AuthToken unmarshal Auth json data.
type AuthToken struct {
	Body struct {
		Token string
	}
}

// MailRuCloud holds all information which required for the api operations.
type MailRuCloud struct {
	Client    *http.Client
	AuthToken string
	Shard     struct {
		Get    string
		Upload string
	}
	IoProgress     func(int64, interface{})
	InitIoProgress func(int64) interface{}
}

// FileList used to unmarshal json information about files in mail.ru cloud.
type FileList struct {
	Email string
	Body  struct {
		Mtime int
		Count struct {
			Folders int
			Files   int
		}
		Tree string
		Name string
		Grev int
		Size int
		Sort struct {
			Order string
			Type  string
		}
		Kind string
		Rev  int
		Type string
		Home string
		Hash string
		List []ListItem
	}
}

// ListItem used to unmarshal json information about a file.
type ListItem struct {
	Name       string
	Size       int
	Kind       string
	Type       string
	Home       string
	Hash       string
	Tree       string
	Mtime      int
	Virus_scan string
	Grev       int
	Rev        int
	Count      struct {
		Folders int
		Files   int
	}
}

// ShardInfo used to unmarshal json information about "shards" which contain urls for api operations.
type ShardInfo struct {
	Email string
	Body  struct {
		Upload []ShardItem
		Get    []ShardItem
	}
	Time   int
	Status int
}

// ShardItem holds information for a particular "shard".
type ShardItem struct {
	Count string
	Url   string
}
