package webmods

import (
  "time"
  "strconv"
  "net/http"
  "io/ioutil"
  "github.com/spf13/viper"
  "github.com/gorilla/mux"
  "github.com/Mstiekema/Yucibot2/base"
)

func UserPage(w http.ResponseWriter, r *http.Request) {
  v := mux.Vars(r)
  usr := v["username"]
  viper.SetConfigFile("../config.toml"); viper.ReadInConfig(); chnl := viper.GetString("twitch.channel")
  
  var db = base.Conn()
  var points, num_lines, userId, timeOffline, timeOnline string
  usrI, _ := db.Query(`SELECT points, num_lines, userId, timeOffline, timeOnline FROM user WHERE name = "`+usr+`"`)
  for usrI.Next() {usrI.Scan(&points, &num_lines, &userId, &timeOffline, &timeOnline)}
  var times, logs []string; var time, log string; lgs, _ := db.Query(`SELECT time, log FROM chatlogs WHERE userId = "`+userId+`" ORDER BY time DESC LIMIT 25`)
  for lgs.Next() {lgs.Scan(&time, &log); times = append(times, time); logs = append(logs, log)}
  r1, _ := http.Get("http://api.yucibot.com/user/pf/"+usr); defer r1.Body.Close(); b1, _ := ioutil.ReadAll(r1.Body); pf := string(b1[:])
  r2, _ := http.Get("http://api.yucibot.com/raw/followsince/"+usr+"/"+chnl); defer r2.Body.Close(); b2, _ := ioutil.ReadAll(r2.Body); fa := string(b2[:])
  r3, _ := http.Get("http://api.yucibot.com/followage/"+usr+"/"+chnl); defer r3.Body.Close(); b3, _ := ioutil.ReadAll(r3.Body); fs := string(b3[:])
  r4, _ := http.Get("http://api.yucibot.com/raw/user/age/"+usr); defer r4.Body.Close(); b4, _ := ioutil.ReadAll(r4.Body); aa := string(b4[:])
  UserInfo := map[string]interface{}{"Username": usr, "Points": points, "Lines": num_lines, "UserId": userId, "TOn": timeOnline, "TOff": timeOffline, "Pf": pf, "AccAge": aa, "FAge": fa, "Fsince": fs, "ChatTimes": times, "ChatLogs": logs}
  db.Close()
  LoadPage(w, r, "./web/templates/user.html", UserInfo)
}

func Logs(w http.ResponseWriter, r *http.Request) {
  var db = base.Conn()
  v := mux.Vars(r)
  t := time.Now(); y := strconv.Itoa(t.Year()); m := strconv.Itoa(int(t.Month())); d := strconv.Itoa(t.Day()); if len(m) == 1 {m = "0"+m}; if len(d) == 1 {d = "0"+d}; date := y+"-"+m+"-"+d
  userId := base.Query("SELECT userId FROM user WHERE name = '"+v["username"]+"'")
  var times, logs []string; var time, log string; lgs, _ := db.Query(`SELECT time, log FROM chatlogs WHERE userId = "`+userId+`" AND DATE(time) = CURDATE() ORDER BY time DESC`)
  for lgs.Next() {lgs.Scan(&time, &log); times = append(times, time); logs = append(logs, log)}
  Logs := map[string]interface{}{"Username": v["username"], "ChatTimes": times, "ChatLogs": logs, "Date": date}
  db.Close()
  LoadPage(w, r, "./web/templates/logs.html", Logs)
}

func LogsDate(w http.ResponseWriter, r *http.Request) {
  var db = base.Conn()
  v := mux.Vars(r)
  var stmt string; 
  date := v["date"]
  if date == "all" {stmt = `ORDER BY time DESC`; date = "All logs"} else {stmt = `AND DATE(time) = "`+date+`" ORDER BY time DESC`}
  userId := base.Query("SELECT userId FROM user WHERE name = '"+v["username"]+"'")
  var times, logs []string; var time, log string; lgs, _ := db.Query(`SELECT time, log FROM chatlogs WHERE userId = "`+userId+`" `+stmt)
  for lgs.Next() {lgs.Scan(&time, &log); times = append(times, time); logs = append(logs, log)}
  Logs := map[string]interface{}{"Username": v["username"], "ChatTimes": times, "ChatLogs": logs, "Date": date}
  db.Close()
  LoadPage(w, r, "./web/templates/logs.html", Logs)
}