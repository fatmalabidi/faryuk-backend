package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"FaRyuk/api/utils"
	"FaRyuk/internal/sharing"
	"FaRyuk/internal/types"
	"FaRyuk/models"

	"github.com/gorilla/mux"
)

func addSharingEndpoints(secure *mux.Router) {
	secure.HandleFunc("/api/share-result", shareResult).Methods("POST")
	secure.HandleFunc("/api/get-pending", getPending).Methods("GET")
	secure.HandleFunc("/api/accept-sharing/{id}", acceptSharing).Methods("GET")
	secure.HandleFunc("/api/decline-sharing/{id}", declineSharing).Methods("GET")
}

func shareResult(w http.ResponseWriter, r *http.Request) {
	var err error

	_, idUser, err := getIdentity(&w, r)
	if err != nil {
		return
	}

	var objmap map[string]json.RawMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteInternalError(&w, unexpectedError)
		return
	}

	err = json.Unmarshal(body, &objmap)
	if err != nil {
		utils.WriteInternalError(&w, "Please provide a valid json")
		return
	}

	var sharedWith string
	err = json.Unmarshal(objmap["idUser"], &sharedWith)
	if err != nil {
		utils.WriteInternalError(&w, "Please provide a 'idUser'")
		return
	}

	var idResult string
	err = json.Unmarshal(objmap["idResult"], &idResult)
	if err != nil {
		utils.WriteInternalError(&w, "Please provide a 'idResult'")
		return
	}
	s := sharing.NewSharing(idUser, idResult, sharedWith)
	dbHandler := models.NewDBHandler()
	defer dbHandler.CloseCommentDBConnection()

	err = dbHandler.InsertSharing(s)

	if err != nil {
		utils.WriteInternalError(&w, dbError)
		return
	}
	utils.ReturnSuccess(&w, "Inserted successfully")
}

func acceptSharing(w http.ResponseWriter, r *http.Request) {
	var err error

	_, idUser, err := getIdentity(&w, r)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	idSharing := vars["id"]

	dbHandler := models.NewDBHandler()
	defer dbHandler.CloseCommentDBConnection()

	s, err := dbHandler.GetSharingByID(idSharing)
	if err != nil {
		utils.WriteInternalError(&w, dbError)
		return
	}

	if s.UserID != idUser {
		utils.WriteForbidden(&w, "Privilege error")
		return
	}

	s.State = "Accepted"
	ok := dbHandler.UpdateSharing(&s)

	if !ok {
		utils.WriteInternalError(&w, dbError)
		return
	}

	res := dbHandler.GetResultByID(s.ResultID)

	if res == nil {
		utils.WriteInternalError(&w, dbError)
		return
	}

	res.SharedWith = append(res.SharedWith, s.UserID)

	ok = dbHandler.UpdateResult(res)

	if !ok {
		utils.WriteInternalError(&w, dbError)
		return
	}
	utils.ReturnSuccess(&w, "Accepted successfully")
}

func declineSharing(w http.ResponseWriter, r *http.Request) {
	var err error

	_, idUser, err := getIdentity(&w, r)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	idSharing := vars["id"]

	dbHandler := models.NewDBHandler()
	defer dbHandler.CloseCommentDBConnection()

	s, err := dbHandler.GetSharingByID(idSharing)
	if err != nil {
		utils.WriteInternalError(&w, dbError)
		return
	}

	if s.UserID != idUser {
		utils.WriteForbidden(&w, "Privilege error")
		return
	}

	s.State = "Declined"
	ok := dbHandler.UpdateSharing(&s)

	if !ok {
		utils.WriteInternalError(&w, dbError)
		return
	}

	utils.ReturnSuccess(&w, "Declined successfully")
}

func getPending(w http.ResponseWriter, r *http.Request) {
	_, idUser, err := getIdentity(&w, r)
	if err != nil {
		return
	}

	dbHandler := models.NewDBHandler()
	defer dbHandler.CloseCommentDBConnection()
	sharings, err := dbHandler.GetSharingsByUser(idUser)

	if err != nil {
		utils.WriteInternalError(&w, dbError)
		return
	}
	results := make([]types.Sharing, 0)
	for idx := range sharings {
		if sharings[idx].State != "Pending" {
			continue
		}
		r := dbHandler.GetResultByID(sharings[idx].ResultID)
		if r == nil {
			continue
		}
		sharings[idx].ResultID = r.Host
		u := dbHandler.GetUserByID(sharings[idx].OwnerID)
		if u == nil {
			continue
		}
		sharings[idx].OwnerID = u.Username
		results = append(results, sharings[idx])
	}

	utils.ReturnSuccess(&w, results)
}
