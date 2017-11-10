package webmods

import (
  "html/template"
  "fmt"
  "net/http"
  "log"
)

func Home(w http.ResponseWriter, r *http.Request){
  if r.URL.Path != "/" {
    t, err := template.New("").ParseFiles("./web/templates/error.html", "./web/templates/header.html")
    if err != nil {
      log.Print("template parsing error: ", err)
    }
    err = t.ExecuteTemplate(w, "base", nil)
    if err != nil {
      log.Print("template executing error: ", err)
    }
    return
  }  
  t, err := template.New("").ParseFiles("./web/templates/home.html", "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", nil)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}

func Test(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "This is a test page xD")
}