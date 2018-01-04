package webmods

import (
  "net/http"
  "github.com/markbates/goth/gothic"
)

func AdminSonglist(w http.ResponseWriter, r *http.Request) {
  LoadAdminPage(w, r, "./web/templates/admin/songlist.html", nil) 
}

func AdminModules(w http.ResponseWriter, r *http.Request){
  LoadAdminPage(w, r, "./web/templates/admin/modules.html", nil)
}

func AdminCommands(w http.ResponseWriter, r *http.Request){
  LoadAdminPage(w, r, "./web/templates/admin/commands.html", nil)
}

func AdminClr(w http.ResponseWriter, r *http.Request) {
  LoadAdminPage(w, r, "./web/templates/admin/clr.html", nil)
}

func LoadAdminPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
  session, _ := gothic.Store.Get(r, "loginSession")
  var lvl int
  
  if session.Values["level"] == nil {
    lvl = 100
  } else {
    lvl = session.Values["level"].(int)
  }
  
  if lvl < 200 {
    LoadPage(w, r, "./web/templates/401.html", nil)
  } else {
    LoadPage(w, r, tmpl, data)
  }
}