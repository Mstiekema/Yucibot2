package webmods

import (
  "os"
  "fmt"
  "net/http"
  "math/rand"
  "io/ioutil"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
  "github.com/markbates/goth/gothic"
  "github.com/Mstiekema/Yucibot2/base"
)

func AdminSonglist(w http.ResponseWriter, r *http.Request) {
  LoadAdminPage(w, r, "./web/templates/admin/songlist.html", nil)
}

func AdminModule(w http.ResponseWriter, r *http.Request){
  db := base.Conn()
  res, err := db.Query("SELECT id, moduleName, shortName, moduleDescription, type, state FROM module WHERE type = 'main'")
  if err != nil {panic(err.Error())}; defer res.Close()
  var ids, modules, shortNames, descs, mTypes, states []string
  for res.Next() {
    var id, module, shortName, desc, mType, state string
    err = res.Scan(&id, &module, &shortName, &desc, &mType, &state)
    ids = append(ids, id)
    modules = append(modules, module)
    shortNames = append(shortNames, shortName)
    descs = append(descs, desc)
    mTypes = append(mTypes, mType)
    states = append(states, state)
  }
  Modules := map[string]interface{}{"id": ids, "module": modules, "shortName": shortNames, "desc": descs, "mType": mTypes, "state": states, "typ": "nil", "values": "nil"}
  db.Close()
  LoadAdminPage(w, r, "./web/templates/admin/modules.html", Modules)
}

func AdminModules(w http.ResponseWriter, r *http.Request){
  db := base.Conn()
  v := mux.Vars(r)
  rows, err := db.Query("SELECT value FROM modulesettings where moduleType = '"+v["type"]+"'"); if err != nil {panic(err.Error())}; defer rows.Close()
  var vals []string
  for rows.Next() {
    var val string
    rows.Scan(&val)
    vals = append(vals, val)
  }
  res, err := db.Query("SELECT id, moduleName, shortName, moduleDescription, type, state FROM module WHERE type = '"+v["type"]+"'"); if err != nil {panic(err.Error())}; defer res.Close()
  var ids, modules, shortNames, descs, mTypes, states []string
  for res.Next() {
    var id, module, shortName, desc, mType, state string
    res.Scan(&id, &module, &shortName, &desc, &mType, &state)
    ids = append(ids, id)
    modules = append(modules, module)
    shortNames = append(shortNames, shortName)
    descs = append(descs, desc)
    mTypes = append(mTypes, mType)
    states = append(states, state)
  }
  Modules := map[string]interface{}{"id": ids, "module": modules, "shortName": shortNames, "desc": descs, "mType": mTypes, "state": states, "typ": v["type"], "values": vals}
  db.Close()
  LoadAdminPage(w, r, "./web/templates/admin/modules.html", Modules)
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

func ModuleAdmin(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)
  if err != nil { http.Error(w, "Something went wrong while trying to alter the module state", http.StatusBadRequest); panic(err) }
  var n map[string]interface{}
  if err := json.Unmarshal([]byte(body), &n); err != nil { fmt.Fprintf(w, "Something went wrong while trying to do your meme request"); panic(err) }

  if n["action"] == "enable" {
    base.Update("module", "state = 1", "moduleName", "'"+n["module"].(string)+"'")
    fmt.Fprintf(w, "Enabled the "+n["module"].(string)+" module");
  } else if n["action"] == "disable" {
    base.Update("module", "state = 0", "moduleName", "'"+n["module"].(string)+"'")
    fmt.Fprintf(w, "Disabled the "+n["module"].(string)+" module");
  } else if n["action"] == "restart" {
    fmt.Fprintf(w, "Shutting down bot.");
    os.Exit(0)
  } else if n["action"] == "update" {
    // This part of the code doesn't work yet.
    base.Update("module", n["id"].(string)+" = "+n["value"].(string), "moduleName", "'"+n["module"].(string)+"'")
    fmt.Fprintf(w, "Updated the "+n["module"].(string)+" module");
  }
}

func PostAdminClr(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)
  if err != nil { http.Error(w, "Something went wrong while trying to do your meme request", http.StatusBadRequest); panic(err) }
  var n map[string]interface{}
  if err := json.Unmarshal([]byte(body), &n); err != nil { fmt.Fprintf(w, "Something went wrong while trying to do your meme request"); panic(err) }

  if n["clrType"] == "meme" {
    db := base.Conn()
    conn, _, err := websocket.DefaultDialer.Dial("ws://"+r.Host+"/post/getCLR/", nil)
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
    conn.Close()
    fmt.Fprintf(w, "Succesfully sent the random meme to the CLR browser");
  } else if n["clrType"] == "forceMeme" {
    db := base.Conn()
    conn, _, err := websocket.DefaultDialer.Dial("ws://"+r.Host+"/post/getCLR/", nil)
    if err != nil { http.Error(w, "Could not open websocket connection", http.StatusBadRequest) }
    meme := base.Query(`SELECT url FROM clr where name = "`+n["name"].(string)+`"`)
    if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "meme", "meme": "`+meme+`"}`)); err != nil {fmt.Println(err);return}
    db.Close()
    conn.Close()
    fmt.Fprintf(w, "Succesfully sent your meme to the CLR browser");
  } else if n["clrType"] == "forceSound" {
    db := base.Conn()
    conn, _, err := websocket.DefaultDialer.Dial("ws://"+r.Host+"/post/getCLR/", nil)
    if err != nil { http.Error(w, "Could not open websocket connection", http.StatusBadRequest) }
    url :=  base.Query(`SELECT url FROM clr where name = "`+n["name"].(string)+`"`)
    if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "sound", "sound": "`+n["name"].(string)+`", "url": "`+url+`"}`)); err != nil {fmt.Println(err);return}
    db.Close()
    conn.Close()
    fmt.Fprintf(w, "Succesfully sent your sound to the CLR browser");
  } else if n["clrType"] == "removeCLR" {
    base.Delete("clr", "id", n["name"].(string))
    fmt.Fprintf(w, "Succesfully removed %s from the database", n["name"]);
  } else if n["clrType"] == "addCLR" {
    base.Insert("clr (name, url, type)", "('"+n["name"].(string)+"', '"+n["url"].(string)+"', '"+n["type"].(string)+"')")
    fmt.Fprintf(w, "Succesfully added %s to the database", n["name"]);
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
