package base

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "fmt"
  "time"
  "github.com/spf13/viper"
)

type AllUsers struct {
  Count int `json:"chatter_count"`
  Chatters struct {
    Mods []string `json:"moderators"`
    Viewers []string `json:"viewers"`
  }
}

type StreamStatus struct {
	Data []struct {
    Title string `json:"title"`
    ViewerCount int `json:"viewer_count"`
		Type string `json:"type"`
	} `json:"data"`
}

func (b *Bot) UpdateUser() {
  var netClient = &http.Client{Timeout: time.Second * 10,}
  viper.SetConfigFile("./config.toml")
  err := viper.ReadInConfig()
  resp, err := netClient.Get("https://tmi.twitch.tv/group/user/"+viper.GetString("twitch.channel")+"/chatters")
  if err != nil {fmt.Println(err); return}
  defer resp.Body.Close()
  
  body, _ := ioutil.ReadAll(resp.Body)
  app := AllUsers{}
  err = json.Unmarshal(body, &app)
  if err != nil {fmt.Println(err); return}
  
  allChatters := make([]string, len(app.Chatters.Mods) + len(app.Chatters.Viewers))
  copy(allChatters, app.Chatters.Mods)
  copy(allChatters[len(app.Chatters.Mods):], app.Chatters.Viewers)
  
  info := b.GetStreamInfo()
  if len(info.Data) != 0 {
    for i := 0; i < len(allChatters); i++ {
      Update("user", "points = points + 5", "name", "'"+allChatters[i]+"'")
      Update("user", "timeOnline = timeOnline + 5", "name", "'"+allChatters[i]+"'")
    }  
  } else {
    for i := 0; i < len(allChatters); i++ {
      Update("user", "timeOffline = timeOffline + 5", "name", "'"+allChatters[i]+"'")
    }
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