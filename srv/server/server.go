package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/rebeccaskinner/agile17-sample/user"
)

// Server represents a server
type Server struct {
	Port     int
	oldUsers map[string]*user.User
	newUsers map[string]*user.NewUser
	oldMutex *sync.Mutex
	newMutex *sync.Mutex
}

// New creates a new server, pre-seeded with data from path (if it's not "")
func New(port int, path string) (*Server, error) {
	srv := &Server{
		Port:     port,
		oldMutex: new(sync.Mutex),
		newMutex: new(sync.Mutex),
		oldUsers: make(map[string]*user.User),
		newUsers: make(map[string]*user.NewUser),
	}
	if path == "" {
		return srv, nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	oldUsers := make([]*user.User, 0)
	if err := json.Unmarshal(data, &oldUsers); err != nil {
		return nil, err
	}
	for _, u := range oldUsers {
		srv.oldUsers[u.ID] = u
	}
	return srv, nil
}

// Run runs the server
func (s *Server) Run() {
	router := httprouter.New()
	router.GET("/oldusers", s.listOldUsers)
	router.GET("/newusers", s.listNewUsers)
	router.GET("/oldusers/:id", s.getOldUser)
	router.GET("/newusers/:id", s.getNewUser)
	router.POST("/oldusers/:id", s.putOldUser)
	router.POST("/newusers/:id", s.putNewUser)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), router))
}

func (s *Server) listOldUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.oldMutex.Lock()
	defer s.oldMutex.Unlock()
	users := make([]string, 0)
	for k := range s.oldUsers {
		users = append(users, k)
	}
	data, err := json.Marshal(users)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
func (s *Server) getOldUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.oldMutex.Lock()
	defer s.oldMutex.Unlock()
	u, ok := s.oldUsers[p.ByName("id")]
	w.Header().Add("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(u)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
func (s *Server) putOldUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.oldMutex.Lock()
	defer s.oldMutex.Unlock()
	uid := p.ByName("id")
	if _, ok := s.oldUsers[uid]; ok {
		w.Write([]byte("user exists"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("error reading http data: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr := &user.User{}
	if err := json.Unmarshal(data, &usr); err != nil {
		logrus.Error("error deserializing json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.oldUsers[uid] = usr
	return
}

func (s *Server) listNewUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.newMutex.Lock()
	defer s.newMutex.Unlock()
	users := make([]string, 0)
	for k := range s.newUsers {
		users = append(users, k)
	}
	data, err := json.Marshal(users)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
func (s *Server) getNewUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.newMutex.Lock()
	defer s.newMutex.Unlock()
	u, ok := s.newUsers[p.ByName("id")]
	w.Header().Add("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(u)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
func (s *Server) putNewUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.newMutex.Lock()
	defer s.newMutex.Unlock()
	uid := p.ByName("id")
	if _, ok := s.newUsers[uid]; ok {
		w.Write([]byte("user exists"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("error reading http data: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr := &user.NewUser{}
	if err := json.Unmarshal(data, &usr); err != nil {
		logrus.Error("error deserializing json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.newUsers[uid] = usr
	return
}
