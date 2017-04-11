package handlers

import (
	"net/http"
	"sort"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rebeccaskinner/agile17-sample/user"
)

// Datastore represents a store of data
type Datastore struct {
	sync.Mutex
	oldData map[string][]*user.User
	newData map[string][]*user.NewUser
}

// NewDatastore loads a new datastore from a directory of json files
func NewDatastore(path string) (*Datastore, error) {
	return nil, nil
}

// DumpOld dumps a list of all oldData users
func (d *Datastore) DumpOld(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

// DumpNew dumps a list of all newData users
func (d *Datastore) DumpNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

// FetchOld gets a single oldUser
func (d *Datastore) FetchOld(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}

// PostNew adds a single user to newData
func (d *Datastore) PostNew(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}

func (d *Datastore) allUsers() []*user.User {
	d.Lock()
	defer d.Unlock()
	allUsers := make([]*user.User, 0)

	// sort the keys in the map before itereating over it to ensure stable
	// ordering for tests.
	keys := []string{}
	for k := range d.oldData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := d.oldData[k]
		allUsers = append(allUsers, v...)
	}
	return allUsers
}
