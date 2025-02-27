package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"FaRyuk/internal/db"
	"FaRyuk/internal/group"
	"FaRyuk/internal/types"

	"github.com/gorilla/mux"
)

func addGroupEndpoints(secure *mux.Router, adminRouter *mux.Router) {
	secure.HandleFunc("/api/get-groups", getGroups).Methods("GET")
	adminRouter.HandleFunc("/api/group", addGroup).Methods("POST")
	adminRouter.HandleFunc("/api/group", deleteGroup).Methods("DELETE")
	adminRouter.HandleFunc("/api/group/user", addUserToGroup).Methods("POST")
	adminRouter.HandleFunc("/api/group/user", removeUserFromGroup).Methods("DELETE")
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	var err error
	var groups []types.Group
	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	username, idUser, err := getIdentity(&w, r)
	if err != nil {
		return
	}

	if username == "admin" {
		groups, err = dbHandler.GetGroups()

		if err != nil {
			writeInternalError(&w, "Database error")
			return
		}
	} else {
		user := dbHandler.GetUserByID(idUser)
		groups = user.Groups
	}

	writeObject(&w, groups)
}

func addGroup(w http.ResponseWriter, r *http.Request) {
	var objmap map[string]json.RawMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInternalError(&w, "Unexpected error")
		return
	}

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		writeInternalError(&w, "Please provide a valid json")
		return
	}

	var name string
	err = json.Unmarshal(objmap["name"], &name)
	if err != nil {
		writeInternalError(&w, "Please provide a 'name'")
		return
	}

	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	groupRes, err := dbHandler.GetGroupsByName(name)
	if err == nil && groupRes.ID != "Dummy" {
		writeInternalError(&w, "Group already exists")
		return
	}

	group := group.NewGroup(name)
	err = dbHandler.InsertGroup(*group)
	if err != nil {
		writeInternalError(&w, "Database error")
		return
	}
	writeObject(&w, "Group added")
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	var objmap map[string]json.RawMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInternalError(&w, "Unexpected error")
		return
	}

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		writeInternalError(&w, "Please provide a valid json")
		return
	}

	var id string
	err = json.Unmarshal(objmap["id"], &id)
	if err != nil {
		writeInternalError(&w, "Please provide a 'id'")
		return
	}

	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	err = dbHandler.RemoveGroupByID(id)
	if err != nil {
		writeInternalError(&w, "Database error")
		return
	}
	writeObject(&w, "Group deleted")
}

func addUserToGroup(w http.ResponseWriter, r *http.Request) {
	var objmap map[string]json.RawMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInternalError(&w, "Unexpected error")
		return
	}

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		writeInternalError(&w, "Please provide a valid json")
		return
	}

	var idGroup string
	err = json.Unmarshal(objmap["idGroup"], &idGroup)
	if err != nil {
		writeInternalError(&w, "Please provide a valid group id")
		return
	}

	var idUser string
	err = json.Unmarshal(objmap["idUser"], &idUser)
	if err != nil {
		writeInternalError(&w, "Please provide a valid user id")
		return
	}

	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	group, err := dbHandler.GetGroupByID(idGroup)
	if err != nil {
		writeInternalError(&w, "Database error - No such group")
		return
	}

	user := dbHandler.GetUserByID(idUser)
	if user == nil {
		writeInternalError(&w, "Database error - No such user")
		return
	}

	for _, grp := range user.Groups {
		if grp.ID == group.ID {
			writeObject(&w, "User already in the group")
			return
		}
	}
	user.Groups = append(user.Groups, group)

	err = dbHandler.UpdateUser(user)
	if err != nil {
		writeInternalError(&w, fmt.Sprintf("%s", err))
		return
	}

	writeObject(&w, "User added to group")
}

func removeUserFromGroup(w http.ResponseWriter, r *http.Request) {
	var objmap map[string]json.RawMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInternalError(&w, "Unexpected error")
		return
	}

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		writeInternalError(&w, "Please provide a valid json")
		return
	}

	var idGroup string
	err = json.Unmarshal(objmap["idGroup"], &idGroup)
	if err != nil {
		writeInternalError(&w, "Please provide a valid group id")
		return
	}

	var idUser string
	err = json.Unmarshal(objmap["idUser"], &idUser)
	if err != nil {
		writeInternalError(&w, "Please provide a valid user id")
		return
	}

	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	user := dbHandler.GetUserByID(idUser)
	if user == nil {
		writeInternalError(&w, "Database error - No such user")
		return
	}

	groups := make([]types.Group, 0)

	for idx := range user.Groups {
		if user.Groups[idx].ID != idGroup {
			groups = append(groups, user.Groups[idx])
		}
	}

	user.Groups = groups
	err = dbHandler.UpdateUser(user)
	if err != nil {
		writeInternalError(&w, fmt.Sprintf("%s", err))
		return
	}

	writeObject(&w, "User removed from group")
}
