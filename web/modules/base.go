package webmods

import (
  "html/template"
  "fmt"
  "net/http"
  "log"
)

func Home(w http.ResponseWriter, r *http.Request){
  // Testing some stuff, some other stuff will be here soon
  channel := "Mstiekema"
  t, err := template.New("").ParseFiles("./web/templates/home.html", "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", channel)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}

func Test(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "This is a test page xD")
}