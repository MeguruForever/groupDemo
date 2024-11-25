package utils

import (
	"GroupDemo/models"
	"sync"
)

var (
	Groups       = make(map[string]string)
	GroupMembers = make(map[string][]string)
	Mu           sync.Mutex
)

var Pool = &models.GroupPools{
	Pools: make(map[string]*models.SingleGroupPool),
}
