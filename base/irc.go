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
  Channel string
  nick string
  C net.Conn
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
    msg := strings.SplitAfterN(mWithUser, ":", 2)[1]
    preUser := strings.SplitAfter(m, ";")[2]
    user := preUser[13:len(preUser)-1]
    fmt.Printf("[CHAT] " + user + ": " + msg)
    
    // Do something with messages 
    // Do something here with tags or smth idk, sub events etc.
    
    // Commands
    var comm string
    i := 1
    if  i >= 1 && i < len(strings.SplitAfter(msg, " ")) {
      comm = strings.SplitN(msg, " ", 2)[0]
    } else {
      comm = strings.TrimSpace(msg)
    }    
    if strings.HasPrefix(msg, "!") == false {return}
    b.Modules(comm, user, msg)
  }
}

func (b *Bot) SendMsg(msg string) {
  fmt.Fprintf(b.C, "PRIVMSG #%s :"+msg+"\r\n", b.Channel,)
  fmt.Println("[SEND] " + msg)
}