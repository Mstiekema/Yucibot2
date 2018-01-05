package webmods

import (
  "fmt"
  "strings"
  "net/http"
  "math/rand"
  "github.com/gorilla/websocket"
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/base"
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
  db := base.Conn()
  res, err := db.Query("SELECT id, name, url, type FROM clr")
  if err != nil {
    panic(err.Error())
  }
  defer res.Close()
  
  var ids []string
  var names []string
  var urls []string
  var types []string
  for res.Next() {
    var id string
    var name string
    var url string
    var sort string
    err = res.Scan(&id, &name, &url, &sort)
    ids = append(ids, id)
    names = append(names, name)
    urls = append(urls, url)
    types = append(types, sort)
  }
  clr := map[string]interface{}{
    "id": ids,
    "name": names,
    "url": urls,
    "type": types,
  }
  db.Close()
  LoadAdminPage(w, r, "./web/templates/admin/clr.html", clr)
}

func PostAdminClr(hub *Hub, w http.ResponseWriter, r *http.Request) {
  read, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
  client := &Client{hub: hub, conn: read, send: make(chan []byte, 256)}
  client.hub.register <- client
  
  if err != nil {
    http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
  }
  for {
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9090/post/getCLR/", nil)
    _, bMsg, err := read.ReadMessage()
    if err != nil { fmt.Println(err); return }
    msg := string(bMsg)
    
    if msg == "meme" {
      db := base.Conn()
      res, err := db.Query(`SELECT url FROM clr where type = "meme"`)
      if err != nil { panic(err.Error()) }
      defer res.Close()
      
      var urls []string
      for res.Next() {
        var url string
        err = res.Scan(&url)
        urls = append(urls, url)
      }
      if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "meme", "meme": "`+urls[rand.Intn(len(urls))]+`"}`)); err != nil {fmt.Println(err);return}
      db.Close()
    } else if strings.Contains(msg, "forceMeme|") {
      db := base.Conn()
      name := strings.SplitAfter(msg, "forceMeme|")[1]
      meme := base.Query(`SELECT url FROM clr where name = "`+name+`"`)
      if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "meme", "meme": "`+meme+`"}`)); err != nil {fmt.Println(err);return}
      db.Close()
    } else if strings.Contains(msg, "forceSound|") {
      db := base.Conn()
      name := strings.SplitAfter(msg, "forceSound|")[1]
      url := base.Query(`SELECT url FROM clr where name = "`+name+`"`)
      if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "sound", "sound": "`+name+`", "url": "`+url+`"}`)); err != nil {fmt.Println(err);return}
      db.Close()
    } else if strings.Contains(msg, "removeCLR|") {
      id := strings.SplitAfter(msg, "removeCLR|")[1]
      base.Delete("clr", "id", id)
    } else if msg == "addSample" {
      // WIP
    }
  }
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