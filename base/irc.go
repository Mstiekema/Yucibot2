package base

import (
  "fmt"
  "bufio"
  "net"
  "strings"
  "github.com/spf13/viper"
)

var channel string

type Bot struct {
  oauth string
  ytApiKey string
  Channel string
  nick string
  C net.Conn
}

type User struct {
  username string
  displayName string
  userId string
  message string
  mod string
  sub string
}

func CrtBot() *Bot {
  viper.SetConfigFile("./config.toml")
  err := viper.ReadInConfig()
  if err != nil {
    fmt.Println(err)
  }
  
  return &Bot{
    oauth: viper.GetString("twitch.oauth"),
    Channel: viper.GetString("twitch.channel"),
    nick: viper.GetString("twitch.botname"),
    ytApiKey: viper.GetString("apiKeys.ytApiKey"),
  }
}

func (b *Bot) Connect() {
  fmt.Println("[DEBUG] Launching bot")  
  conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
  if err != nil {
  	fmt.Println(err)
  }
  b.C = conn
  
  fmt.Fprintf(conn, "PASS %s\r\n", b.oauth)
  fmt.Fprintf(conn, "NICK %s\r\n", b.nick)
  fmt.Fprintf(conn, "JOIN #%s\r\n", b.Channel)
  fmt.Fprintf(conn, "CAP REQ :twitch.tv/membership\r\n")
  fmt.Fprintf(conn, "CAP REQ :twitch.tv/tags\r\n")
  fmt.Fprintf(conn, "CAP REQ :twitch.tv/commands\r\n")
  fmt.Println("[DEBUG] Connected to chat")
  fmt.Fprintf(conn, "PRIVMSG #%s :Hey there, Yucibot 2.0 is now online KKona\r\n", b.Channel)
}

func (b *Bot) Reader(conn net.Conn) {
  reader := bufio.NewReader(conn)
  for {
    line, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println("err")
      fmt.Println(err)
      break
    }
    if strings.HasPrefix(line, "PING") {
      fmt.Fprintf(conn, "PONG \r\n")
    }
    b.parseMsg(line)
  }
}

func (b *Bot) parseMsg(m string) {
  i := 1
  if i >= 1 && i < len(strings.SplitAfter(m, "PRIVMSG")) {
    mWithUser := strings.SplitAfter(m, "PRIVMSG")[1]
    mWithUser = strings.TrimPrefix(mWithUser, " #")
    msg := strings.TrimSpace(strings.SplitAfterN(mWithUser, ":", 2)[1])
    user := strings.SplitAfter(m, ";")[2][13:len(strings.SplitAfter(m, ";")[2])-1]
    fmt.Println("[CHAT] " + user + ": " + msg)
    
    User := User{}
    User.username = strings.ToLower(user)
    User.displayName = user
    User.userId = strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "user-id=")[1], ";")[0], ";")
    User.message = msg
    User.mod = strings.SplitAfter(m, ";")[5][4:len(strings.SplitAfter(m, ";")[5])-1]
    User.sub = strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "subscriber=")[1], ";")[0], ";")
    
    // Events
    b.UpdateLines(User)
    b.Links(User)
    
    // Commands
    var comm string
    i := 1
    if  i >= 1 && i < len(strings.SplitAfter(msg, " ")) {
      comm = strings.SplitN(msg, " ", 2)[0]
    } else {
      comm = msg
    }    
    if strings.HasPrefix(msg, "!") == false {return}
    b.Modules(comm, User)
  }
}

func (b *Bot) SendMsg(msg string) {
  fmt.Fprintf(b.C, "PRIVMSG #%s :"+msg+"\r\n", b.Channel,)
  fmt.Println("[SEND] " + msg)
}

func (b *Bot) SendWhisper(msg, user string) {
  fmt.Fprintf(b.C, "PRIVMSG #%s : .w "+user+" "+msg+"\r\n", b.Channel,)
  fmt.Println("[WHISPER] "+user+": "+msg)
}