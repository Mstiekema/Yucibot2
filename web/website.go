package web

import (
  "strings"
  "html/template"
  "time"
  "strconv"
  "fmt"
  "net/http"
  "log"
  "github.com/Mstiekema/Yucibot2/base"
)

type User struct {
	Username string
  Points string
  Lines string
}

func MainWeb() {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
  
  http.HandleFunc("/", home)
  http.HandleFunc("/test", test)
  http.HandleFunc("/user/", userPage)
  http.HandleFunc("/songlist/", songlist)
  
  err := http.ListenAndServe(":9090", nil)
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}

func home(w http.ResponseWriter, r *http.Request){
  // Testing some stuff, some other stuff will be here soon
  channel := "Mstiekema"
  t, err := template.ParseFiles("./web/templates/home.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.Execute(w, channel)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}

func userPage(w http.ResponseWriter, r *http.Request) {
  usr := strings.Replace(strings.SplitAfter(r.URL.Path, "/")[2], "/", "", 2)
  points := base.Query("user", "points", usr)
  lines := base.Query("user", "num_lines", usr)
  u := &User{Username: usr, Points: points, Lines: lines}
  
  t, err := template.ParseFiles("./web/templates/user.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.Execute(w, u)
  if err != nil {
    log.Print("template executing error: ", err)
  }
}

func songlist(w http.ResponseWriter, r *http.Request) {
  date := strings.Replace(strings.SplitAfter(r.URL.Path, "/")[2], "/", "", 2)
  if date == "" {
    y, m, d := time.Now().Date()
    date = strconv.Itoa(y)+"-"+strconv.Itoa(int(m))+"-"+strconv.Itoa(d)
  }
  
  var db = base.Conn()
  res, err := db.Query("SELECT name, title, thumb, length, songid, playState FROM songrequest WHERE DATE(time) = '"+date+"'")
  if err != nil {
    panic(err.Error())
  }
  defer res.Close()
  
  //Songlist
  var names []string
  var titles []string
  var thumbs []string
  var lengths []string
  var songids []string
  var playStates []string
  var currSongT string
  var currSongId string
  var currSongN string
  
  for res.Next() {
    var name string
    var title string
    var thumb string
    var length string
    var songid string
    var playState string
    err = res.Scan(&name, &title, &thumb, &length, &songid, &playState)
    names = append(names, name)
    titles = append(titles, title)
    thumbs = append(thumbs, thumb)
    lengths = append(lengths, length)
    songids = append(songids, songid)
    playStates = append(playStates, playState)
  }
  
  for i, v := range playStates {
    if v == "0" {
      currSongT = titles[i]
      currSongId = songids[i]
      currSongN = names[i]
      break
    } else {
      currSongT = ""
      currSongId = ""
      currSongN = ""
    }
  }
  
  Songs := map[string]interface{}{
    "Date": date,
    "CurrSongN": currSongN,
    "CurrSongT": currSongT,
    "CurrSongId": currSongId,
    "Name": names,
    "Title": titles,
    "Thumb": thumbs,
    "Length": lengths,
    "Songid": songids,
    "PlayState": playStates,
  }
  
  t, err := template.ParseFiles("./web/templates/songlist.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.Execute(w, Songs)
  if err != nil {
    log.Print("template executing error: ", err)
  }
  db.Close()
}

func test(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "This is a test page xD")
}