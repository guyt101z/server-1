package api

import (
	"encoding/json"
	"net/http"
	"log"
	"fmt"
	"strconv"
	"database/sql"
	"github.com/OPENCBS/server/iface"
	"github.com/OPENCBS/server/model"
	"github.com/OPENCBS/server/util"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var js []byte
	var repo iface.UserRepo
	var items []*model.User
	var users *model.Users
	var offset int
	var limit int

	db, err := iface.GetDb(r)
	if err != nil {
		goto Error
	}

	repo = iface.NewUserRepo()
	offset = util.GetOffset(r)
	limit = util.GetLimit(r)
	items, err = repo.FindAll(db, offset, limit)
	if err != nil {
		goto Error
	}
	for _, user := range items {
		user.Href = fmt.Sprintf("%s/users/%d", util.GetBaseUrl(r), user.Id)
	}
	users = new(model.Users)
	users.Href = fmt.Sprintf("%s/users", util.GetBaseUrl(r))
	users.Offset = offset
	users.Limit = limit
	users.Items = items

	js, err = json.MarshalIndent(users, "", "  ")
	if err != nil {
		goto Error
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(js)
	return

Error:
	log.Println(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var js []byte
	var repo iface.UserRepo
	var user *model.User
	var db *sql.DB

	idString := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		goto Error
	}

	db, err = iface.GetDb(r)
	if err != nil {
		goto Error
	}

	repo = iface.NewUserRepo()
	user, err = repo.FindById(db, id)
	if err != nil {
		goto Error
	}
	if user == nil {
		goto NotFound
	}
	user.Href = fmt.Sprintf("%s/users/%d", util.GetBaseUrl(r), user.Id)

	js, err = json.MarshalIndent(user, "", "  ")
	if err != nil {
		goto Error
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(js)
	return

NotFound:
	http.Error(w, "User not found", http.StatusNotFound)
	return

Error:
	log.Println(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

