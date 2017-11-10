package webmods

import (
  "strings"
  "html/template"
  "net/http"
  "log"
  "github.com/Mstiekema/Yucibot2/base"
)

type User struct {
	Username string
  Points string
  Lines string
}

func UserPage(w http.ResponseWriter, r *http.Request) {
  usr := strings.Replace(strings.SplitAfter(r.URL.Path, "/")[2], "/", "", 2)
  points := base.Query("user", "points", usr)
  lines := base.Query("user", "num_lines", usr)
  u := &User{Username: usr, Points: points, Lines: lines}
  
  t, err := template.New("").ParseFiles("./web/templates/user.html", "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", u)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}