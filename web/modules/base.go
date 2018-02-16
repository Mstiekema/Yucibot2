package webmods

import (
  "log"
  "math"
  "strconv"
  "net/http"
  "html/template"
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/base"
)

func Home(w http.ResponseWriter, r *http.Request){
  LoadPage(w, r, "./web/templates/home.html", nil)
}

func Commands(w http.ResponseWriter, r *http.Request){
  db := base.Conn()
  res, err := db.Query("SELECT level, commName, COALESCE(commDesc, '') as commDesc, COALESCE(commUse, '') as commUse, COALESCE(response, '') as response, points, cd FROM commands")
  if err != nil {
    panic(err.Error())
  }
  defer res.Close()
  
  var levels, commNames, commDescs, commUses, responses, pointss, cds []string
  for res.Next() {
    var level, commName, commDesc, commUse, response, points, cd string
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
  db.Close()
  LoadPage(w, r, "./web/templates/commands.html", commands)
}

func Stats(w http.ResponseWriter, r *http.Request){
  db := base.Conn()
  var lines, userCount, songRequests, timeoutCount, banCount string
  
  l, _ := db.Query(`SELECT COUNT(*) FROM chatlogs`); for l.Next() {l.Scan(&lines)}
  u, _ := db.Query(`SELECT COUNT(*) FROM user`); for u.Next() {u.Scan(&userCount)}
  s, _ := db.Query(`SELECT COUNT(*) FROM songrequest`); for s.Next() {s.Scan(&songRequests)}
  t, _ := db.Query(`SELECT COUNT(*) FROM adminlogs where type = "timeout"`); for t.Next() {t.Scan(&timeoutCount)}
  b, _ := db.Query(`SELECT COUNT(*) FROM adminlogs where type = "ban"`); for b.Next() {b.Scan(&banCount)}
  
  var nlNames, totalLines, pNames, totalPoints []string
  // onlineHours, onNames, offlineHours, offNames []string
  var nlName, totalLine, pName, totalPoint string
  // onlineHour, onName, offlineHour, offName string
  tl, _ := db.Query("SELECT name, num_lines FROM user ORDER BY num_lines DESC LIMIT 15")
  for tl.Next() {tl.Scan(&nlName, &totalLine); totalLines = append(totalLines, totalLine); nlNames = append(nlNames, nlName)}
  tp, _ := db.Query("SELECT name, points FROM user ORDER BY points DESC LIMIT 15")
  for tp.Next() {tp.Scan(&pName, &totalPoint); totalPoints = append(totalPoints, totalPoint); pNames = append(pNames, pName)}
  // ton, _ := db.Query("SELECT name, timeOnline FROM user ORDER BY timeOnline DESC LIMIT 15")
  // for ton.Next() {ton.Scan(&onName, &onlineHour); onlineHours = append(onlineHours, onlineHour); onNames = append(onNames, onName)}
  // tof, _ := db.Query("SELECT name, timeOffline FROM user ORDER BY timeOffline DESC LIMIT 15")
  // for tof.Next() {tof.Scan(&offName, &offlineHour); offlineHours = append(offlineHours, offlineHour); offNames = append(offNames, offName)}
  
  Stats := map[string]interface{}{
    "lines": lines,
    "userCount": userCount,
    "songRequests": songRequests,
    "timeoutCount": timeoutCount,
    "banCount": banCount,
    "totalLines": totalLines,
    "nlNames": nlNames,
    "totalPoints": totalPoints,
    "pNames": pNames,
    // "onlineHours": onlineHours,
    // "onNames": onNames,
    // "offlineHours": offlineHours,
    // "offNames": offNames,
  }
  db.Close()
  LoadPage(w, r, "./web/templates/stats.html", Stats)
}

func Error(w http.ResponseWriter, r *http.Request) {
  LoadPage(w, r, "./web/templates/404.html", nil)
}

func LoadPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
  parseWatchTime := func (y string) string {
    // Fix deze shit, want negative minutes klopt niet helemaal he
    x, _ := strconv.Atoi(y); days := x / 60 / 24; hours := 24 + int(math.Remainder(float64(x / 60), 24)); minutes := x - days*24*60 - hours*60
    return strconv.Itoa(days)+" days, "+strconv.Itoa(hours)+" hours, "+strconv.Itoa(minutes)+" minutes"
  }
  add := func (x, y int) int { return x + y }
  funcs := template.FuncMap{"add": add, "parseWatchTime": parseWatchTime}
  session, _ := gothic.Store.Get(r, "loginSession")
  if data != nil {
    session.Values["Info"] = data
  }
  t, err := template.New("").Funcs(funcs).ParseFiles(tmpl, "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", session)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}