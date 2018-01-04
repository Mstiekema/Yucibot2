package webmods

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/Mstiekema/Yucibot2/base"
)

type User struct {
	Username string
  Points string
  Lines string
}

func UserPage(w http.ResponseWriter, r *http.Request) {
  v := mux.Vars(r)
  usr := v["username"]
  points := base.Query("SELECT points FROM user WHERE name = '"+usr+"'")
  lines := base.Query("SELECT num_lines FROM user WHERE name = '"+usr+"'")
  u := &User{Username: usr, Points: points, Lines: lines}
  LoadPage(w, r, "./web/templates/user.html", u)
}