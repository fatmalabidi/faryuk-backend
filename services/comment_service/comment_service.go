package comment_service

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"FaRyuk/internal/comment"
	"FaRyuk/internal/db"
	"FaRyuk/internal/group"
	"FaRyuk/internal/helper"
	"FaRyuk/internal/types"
	"FaRyuk/internal/user"
	"FaRyuk/services/comment_service/commun"
	"FaRyuk/utils"

	"github.com/gorilla/mux"
)

const adminUsername = "admin"

func GetComments(w http.ResponseWriter, r *http.Request) {
  var comments []types.Comment
  var result *types.Result
  var err error
  vars := mux.Vars(r)
  idResult := vars["id"]

  token, err := utils.GetCookie("jwt-token", r)
  if err != nil {
    utils.WriteForbidden(&w, "Cookie error")
    return
  }

  username, _, err := user.GetUsername(token, commun.JWTSecret)
  if err != nil {
    utils.WriteForbidden(&w, "Token error")
    return
  }

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  if username == adminUsername {
    comments, err = dbHandler.GetCommentsByResult(idResult)
  } else {
    result = dbHandler.GetResultByID(idResult)
    if result == nil {
      err = fmt.Errorf("cannot retreive comments for this result")
    } else {
      comments, err = dbHandler.GetCommentsByResult(idResult)
    }
  }

  if err != nil {
    utils.WriteInternalError(&w, "Database error")
    return
  }
  utils.WriteObject(&w, comments)
}

func InsertComment(w http.ResponseWriter, r *http.Request) {
  var objmap map[string]json.RawMessage
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    utils.WriteInternalError(&w, "Unexpected error")
    return
  }

  err = json.Unmarshal(body, &objmap)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide a valid json")
    return
  }

  var content string
  err = json.Unmarshal(objmap["content"], &content)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide a 'content'")
    return
  }

  var idOwner string
  err = json.Unmarshal(objmap["owner"], &idOwner)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide an 'owner'")
    return
  }

  var idResult string
  err = json.Unmarshal(objmap["idResult"], &idResult)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide a 'result'")
    return
  }

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  if content[0] == '#' {
    content = strings.ToLower(content)
    err = dbHandler.AddTagsToResult(idResult, helper.GetTags(content))
  } else {
    comment := comment.NewComment(html.EscapeString(content), idOwner, idResult)
    err = dbHandler.InsertComment(comment)
  }

  if err != nil {
    utils.WriteInternalError(&w, "Database error")
    return
  }
  utils.WriteObject(&w, "Comment posted")
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
  var objmap map[string]json.RawMessage
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    utils.WriteInternalError(&w, "Unexpected error")
    return
  }

  err = json.Unmarshal(body, &objmap)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide a valid json")
    return
  }

  var idComment string
  err = json.Unmarshal(objmap["id"], &idComment)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide an 'id'")
    return
  }

  var content string
  err = json.Unmarshal(objmap["content"], &content)
  if err != nil {
    utils.WriteInternalError(&w, "Please provide a 'content'")
    return
  }

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  comment, err := dbHandler.GetCommentByID(idComment)
  if err != nil {
    utils.WriteNotFound(&w, "Comment not found")
    return
  }

  comment.Content = html.EscapeString(content)
  comment.UpdatedDate = time.Now()
  ok := dbHandler.UpdateComment(&comment)

  if !ok {
    utils.WriteInternalError(&w, "Could not update record")
    return
  }
  utils.WriteObject(&w, "Updated successfully")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  res := dbHandler.RemoveCommentByID(id)
  if !res {
    utils.WriteInternalError(&w, "Could not delete record")
    return
  }
  utils.WriteObject(&w, "Deleted successfully")
}

func GetCommentsHighlight(w http.ResponseWriter, r *http.Request) {
  var results []types.Result
  var err error
  pageSize := 10
  offset := 1
  search := ""

  requestedComments := make([]types.Comment,0)
  resultsIds := make([]string, 0)
  query := r.URL.Query()

  // Check user's permission
  token, err := utils.GetCookie("jwt-token", r)
  if err != nil {
    utils.WriteForbidden(&w, "Cookie error")
    return
  }

  username, idUser, err := user.GetUsername(token, commun.JWTSecret)
  if err != nil {
    utils.WriteForbidden(&w, "Token error")
    return
  }

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  // Get results of current user
  if username == adminUsername {
    results, err = dbHandler.GetResultsBySearch(make(map[string]string), -1, -1)
  } else {
    user := dbHandler.GetUserByID(idUser)
    results, err = dbHandler.GetResultsBySearchAndOwner(make(map[string]string), idUser, group.ToIDsArray(user.Groups), -1, -1)
  }
  if err != nil {
    utils.WriteInternalError(&w, "Database error")
    return
  }

  // Put result ids in slice
  for _, r := range results {
    resultsIds = append(resultsIds, r.ID)
  }

  searchSlice, ok := query["search"]
  if ok {
    search = searchSlice[0]
  }

  // Get all comments by chronological order
  comments, err := dbHandler.GetCommentsByText(search)
  if err != nil {
    utils.WriteInternalError(&w, "Database error")
  }

  helper.Reverse(comments)

  // Check how many the user asked for
  pageSizeSlice, ok := query["size"]
  if ok {
    pageSize, _ = strconv.Atoi(pageSizeSlice[0])
  }

  offsetSlice, ok := query["offset"]
  if ok {
    offset, _ = strconv.Atoi(offsetSlice[0])
  }
  taken := 0
  for _, c := range comments {
    if len(requestedComments) >= pageSize {
      break
    }
    if !helper.ContainsStr(resultsIds, c.IDResult) {
      continue
    }
    taken++
    if taken < offset {
      continue
    }
    requestedComments = append(requestedComments, c)
  }
  utils.WriteObject(&w, requestedComments)
}

func GetCountHighlight(w http.ResponseWriter, r *http.Request) {
  var results []types.Result
  var err error

  resultsIds := make([]string, 0)
  query := r.URL.Query()

  // Check user's permission
  token, err := utils.GetCookie("jwt-token", r)
  if err != nil {
    utils.WriteForbidden(&w, "Cookie error")
    return
  }

  username, idUser, err := user.GetUsername(token, commun.JWTSecret)
  if err != nil {
    utils.WriteForbidden(&w, "Token error")
    return
  }

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  // Get results of current user
  if username == adminUsername {
    results, err = dbHandler.GetResultsBySearch(make(map[string]string), -1, -1)
  } else {
    user := dbHandler.GetUserByID(idUser)
    results, err = dbHandler.GetResultsBySearchAndOwner(make(map[string]string), idUser, group.ToIDsArray(user.Groups), -1, -1)
  }
  if err != nil {
    utils.WriteInternalError(&w, "Database error")
    return
  }

  // Put result ids in slice
  for _, r := range results {
    resultsIds = append(resultsIds, r.ID)
  }

  searchSlice, ok := query["search"]
  search := ""
  if ok {
    search = searchSlice[0]
  }

  // Get all comments by chronological order
  comments, err := dbHandler.GetCommentsByText(search)
  if err != nil {
    utils.WriteInternalError(&w, "Database error")
  }

  taken := 0
  for _, c := range comments {
    if helper.ContainsStr(resultsIds, c.IDResult) {
      taken++
    }
  }
  utils.WriteObject(&w, taken)
}

func GetTags(w http.ResponseWriter, r *http.Request) {
  var results []types.Result
  var err error
  tags := make(map[string]int)

  // Check user's permission
  token, err := utils.GetCookie("jwt-token", r)
  if err != nil {
    utils.WriteForbidden(&w, "Cookie error")
    return
  }

  username, idUser, err := user.GetUsername(token, commun.JWTSecret)
  if err != nil {
    utils.WriteForbidden(&w, "Token error")
    return
  }

  dbHandler := db.NewDBHandler()
  defer dbHandler.CloseConnection()

  // Get results of current user
  if username == adminUsername {
    results, err = dbHandler.GetResultsBySearch(make(map[string]string), -1, -1)
  } else {
    user := dbHandler.GetUserByID(idUser)
    results, err = dbHandler.GetResultsBySearchAndOwner(make(map[string]string), idUser, group.ToIDsArray(user.Groups), -1, -1)
  }
  if err != nil {
    utils.WriteInternalError(&w, "Database error")
    return
  }

  // Put result ids in slice
  for _, r := range results {
    for _, tag := range r.Tags {
      tags[tag]++
    }
  }

  utils.WriteObject(&w, tags)
}
