package webmods

import (
  "strings"
  "time"
  "fmt"
  "strconv"
  "encoding/json"
  "net/http"
  
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/base"
)

func getSongs(qry, date string) map[string]interface {} {
  var db = base.Conn()
  res, err := db.Query(qry+"'"+date+"'")
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
    "Msg": nil,
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
  
  db.Close()
  return Songs
}

func Songlist(w http.ResponseWriter, r *http.Request) {
  v := mux.Vars(r)
  date := v["date"]
  var Songs = getSongs("SELECT name, title, thumb, length, songid, playState FROM songrequest WHERE DATE(time) =", date)
  LoadPage(w, r, "./web/templates/songlist.html", Songs) 
}

func TodaySonglist(w http.ResponseWriter, r *http.Request) {
  y, m, d := time.Now().Date()
  date := strconv.Itoa(y)+"-"+strconv.Itoa(int(m))+"-"+strconv.Itoa(d)
  w.Header().Set("Location", "/songlist/"+date)
  w.WriteHeader(http.StatusTemporaryRedirect)
}

func AdminSonglist(w http.ResponseWriter, r *http.Request) {
  session, _ := gothic.Store.Get(r, "loginSession")
  var lvl int
  
  if session.Values["level"] == nil {
    lvl = 100
  } else {
    lvl, _ = strconv.Atoi(session.Values["level"].(string))
  }
  
  if lvl < 200 {
    LoadPage(w, r, "./web/templates/401.html", nil)
    return
  }
  
  LoadPage(w, r, "./web/templates/adminSonglist.html", nil) 
}

func SendSongs(w http.ResponseWriter, r *http.Request) {
  conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
  if err != nil {
    http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
  }
  for {
    // Get message (_ is msg)
    p, bMsg, err := conn.ReadMessage()
    if err != nil {
      fmt.Println(err)
      return
    }
    msg := string(bMsg)
    
    // Replies
    if (msg == "refreshData") {
      y, m, d := time.Now().Date()
      date := strconv.Itoa(y)+"-"+strconv.Itoa(int(m))+"-"+strconv.Itoa(d)
      Songs := getSongs("SELECT name, title, thumb, length, songid, playState FROM songrequest WHERE playState = 0 AND DATE(time) =", date)
      Songs["Msg"] = "pushSonglist"
      jsonSongs, _ := json.Marshal(Songs)

      if err := conn.WriteMessage(p, jsonSongs); err != nil {
        fmt.Println(err)
        return
      }
    } else if (msg == "prevSong") {
      var db = base.Conn()
      y, m, d := time.Now().Date()
      date := strconv.Itoa(y)+"-"+strconv.Itoa(int(m))+"-"+strconv.Itoa(d)
      res, err := db.Query("select songid from songrequest where playState = 1 AND DATE(time) = '"+date+"' ORDER BY id DESC LIMIT 1")
      if err != nil {
        panic(err.Error())
      }
      defer res.Close()
      
      var oldId string
      for res.Next() {
        var songid string
        err = res.Scan(&songid)
        oldId = songid
      }
      db.Close()
      
      base.Update("songrequest", "playState = 0", "songid", "'"+oldId+"'")
      Songs := getSongs("SELECT name, title, thumb, length, songid, playState FROM songrequest WHERE playState = 0 AND DATE(time) =", date)
      Songs["Msg"] = "prevSongInfo"
      jsonSongs, _ := json.Marshal(Songs)

      if err := conn.WriteMessage(p, jsonSongs); err != nil {
        fmt.Println(err)
        return
      }
    } else if (strings.Contains(msg, "endSong|")) {
      id := strings.SplitAfter(msg, "endSong|")[1]
      base.Update("songrequest", "playState = 1", "songid", "'"+id+"'")
      state, _ := json.Marshal("'nextSong'")
      if err := conn.WriteMessage(p, state); err != nil {
        fmt.Println(err)
        return
      }
    } else if (strings.Contains(msg, "delSong|")) {
      id := strings.SplitAfter(msg, "delSong|")[1]
      base.Update("songrequest", "playState = 1", "songid", "'"+id+"'")
      state, _ := json.Marshal("'confDelSong'")
      if err := conn.WriteMessage(p, state); err != nil {
        fmt.Println(err)
        return
      }
    }
  }
}