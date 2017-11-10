package webmods

import (
  "strings"
  "html/template"
  "time"
  "strconv"
  "net/http"
  "log"
  "github.com/Mstiekema/Yucibot2/base"
)

func Songlist(w http.ResponseWriter, r *http.Request) {
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
  
  t, err := template.New("").ParseFiles("./web/templates/songlist.html", "./web/templates/header.html")
  if err != nil {
    log.Print("template parsing error: ", err)
  }
  err = t.ExecuteTemplate(w, "base", Songs)
  if err != nil {
    log.Print("template executing error: ", err)
  }
  db.Close()
}