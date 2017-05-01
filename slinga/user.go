package slinga

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

/*
 	This file declares all the necessary structures for Users to be retrieved
 	For now it loads users with labels from a YAML file
 	Later this will be replaced with LDAP integration
  */

type User struct {
	Id       string
	Name     string
	Labels	 LabelSet
}

type GlobalUsers struct {
	Users map[string]User
}

// Loads users from YAML file
func loadUserByIDFromDir(dir string, id string) User {
	return loadUsersFromDir(dir).Users[id]
}

// Loads users from YAML file
func loadUsersFromDir(dir string) GlobalUsers {
	dat, e := ioutil.ReadFile(dir + "users.yaml")
	if e != nil {
		panic(e)
	}
	t := []User{}
	e = yaml.Unmarshal([]byte(dat), &t)
	if e != nil {
		log.Fatalf("error: %v", e)
	}
	r := GlobalUsers{Users: make(map[string]User)}
	for _, u := range t {
		r.Users[u.Id] = u;
	}
	return r
}