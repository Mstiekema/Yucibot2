package base

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "fmt"
  "time"
)

type AllUsers struct {
  Count int `json:"chatter_count"`
  Chatters struct {
    Mods []string `json:"moderators"`
    Viewers []string `json:"viewers"`
  }
}

func (b *Bot) UpdatePoints() {
  var netClient = &http.Client{Timeout: time.Second * 10,}
  resp, _ := netClient.Get("https://tmi.twitch.tv/group/user/merijn/chatters")
  defer resp.Body.Close()
  
  body, _ := ioutil.ReadAll(resp.Body)
  app := AllUsers{}
  err := json.Unmarshal(body, &app)
  if err != nil {
    fmt.Println(err)
    return
  }
  
  allChatters := make([]string, len(app.Chatters.Mods) + len(app.Chatters.Viewers))
  copy(allChatters, app.Chatters.Mods)
  copy(allChatters[len(app.Chatters.Mods):], app.Chatters.Viewers)
  
  for i := 0; i < len(allChatters); i++ {
    Update("user", "points = points + 5", "name", "'"+allChatters[i]+"'")
  }
}

func (b *Bot) UpdateLines(U User) {
  StringOldPoints := Query("SELECT points FROM user WHERE name = '"+U.username+"'")
  if StringOldPoints != "" {
    Update("user", "num_lines = num_lines + 1", "name", "'"+U.username+"'")
    Insert("chatlogs (userId, log)", "('"+U.userId+"', '" +U.message+"')")
  } else {
    Insert("user (name, userId, points, num_lines, level, isMod)", "('"+U.username+"', '"+U.userId+"', '0', '1', '100', '"+U.mod+"')")
    Insert("chatlogs (userId, log)", "('"+U.userId+"', '" +U.message+"')")
  }
}