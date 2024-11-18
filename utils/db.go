package utils

import "sync"

var (
	Groups       = make(map[string]string)
	GroupMembers = make(map[string][]string)
	Messages     = make(map[string][]string)
	Mu           sync.Mutex
)
