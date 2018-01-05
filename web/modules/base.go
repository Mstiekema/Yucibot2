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

func Commands(w http.ResponseWriter, r *http.Request){
  res, err := db.Query("SELECT level, commName, COALESCE(commDesc, '') as commDesc, COALESCE(commUse, '') as commUse, COALESCE(response, '') as response, points, cd FROM commands")
  if err != nil {
    panic(err.Error())
  }
  defer res.Close()
  
  var levels []string
  var commNames []string
  var commDescs []string
  var commUses []string
  var responses []string
  var pointss []string
  var cds []string
  for res.Next() {
    var level string
    var commName string
    var commDesc string
    var commUse string
    var response string
    var points string
    var cd string
    err = res.Scan(&level, &commName, &commDesc, &commUse, &response, &points, &cd)
    levels = append(levels, level)
    commNames = append(commNames, commName)
    commDescs = append(commDescs, commDesc)
    commUses = append(commUses, commUse)
    responses = append(responses, response)
    pointss = append(pointss, points)
    cds = append(cds, cd)
  }
  commands := map[string]interface{}{
    "level": levels,
    "commName": commNames,
    "commDesc": commDescs,
    "commUse": commUses,
    "response": responses,
    "points": pointss,
    "cd": cds,
  }
  
  LoadPage(w, r, "./web/templates/commands.html", commands)
}

func Stats(w http.ResponseWriter, r *http.Request){
  LoadPage(w, r, "./web/templates/stats.html", nil)
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