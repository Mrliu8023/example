package search

import groupsdev "example/groups-dev"

type Searcher interface {
	Search(keyword string) (*groupsdev.GroupList, error)
	AddDocuments(docs []Document) error
}

type Document interface {
	Name() string
	Content() string
}
