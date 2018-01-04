package webmods

import (
  "html/template"
  "net/http"
  "log"
  "github.com/markbates/goth/gothic"
)

func Home(w http.ResponseWriter, r *http.Request){
  LoadPage(w, r, "./web/templates/home.html", nil)
}

func Error(w http.ResponseWriter, r *http.Request) {
  LoadPage(w, r, "./web/templates/404.html", nil)
}

func LoadPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
  session, _ := gothic.Store.Get(r, "loginSession")
  if data != nil {
    session.Values["Info"] = data
  }
  t, err := template.New("").ParseFiles(tmpl, "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", session)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}