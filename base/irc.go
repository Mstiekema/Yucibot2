package base

import (
  "os"
  "fmt"
  "net"
  "bufio"
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
      os.Exit(0)
    }
    if strings.HasPrefix(line, "PING") {
      fmt.Fprintf(conn, "PONG \r\n")
    }
    b.ParseMsg(line)
  }
}

func (b *Bot) ParseMsg(m string) {
  if 1 < len(strings.SplitAfter(m, "PRIVMSG")) {
    mWithUser := strings.SplitAfter(m, "PRIVMSG")[1]
    mWithUser = strings.TrimPrefix(mWithUser, " #")
    msg := strings.TrimSpace(strings.SplitAfterN(mWithUser, ":", 2)[1])
    user := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "display-name=")[1], ";")[0], ";")
    fmt.Println("[CHAT] " + user + ": " + msg)

    User := User{}
    User.username = strings.ToLower(user)
    User.displayName = user
    User.userId = strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "user-id=")[1], ";")[0], ";")
    User.message = msg
    User.mod = strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "mod=")[1], ";")[0], ";")
    User.sub = strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "subscriber=")[1], ";")[0], ";")

    var comm string
    i := 1
    if  i >= 1 && i < len(strings.SplitAfter(msg, " ")) {
      comm = strings.SplitN(msg, " ", 2)[0]
    } else {
      comm = msg
    }

    b.UpdateLines(User)
    b.Links(User)
    b.Nuke(comm, User, msg)

    if strings.HasPrefix(msg, "!") == false {return}
    b.Modules(comm, User)
  }
  if 1 < len(strings.SplitAfter(m, "USERNOTICE")) {
    noticeType := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-id=")[1], ";")[0], ";")
    if noticeType == "sub" || noticeType == "resub" {
      username := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "login=")[1], ";")[0], ";")
      displayName := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "display-name=")[1], ";")[0], ";")
      months := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-param-months=")[1], ";")[0], ";")
      plan := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-param-sub-plan=")[1], ";")[0], ";")
      b.Subs(noticeType, username, displayName, months, plan)
    } else if noticeType == "subgift" {
      giftU := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "login=")[1], ";")[0], ";")
      giftD := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "display-name=")[1], ";")[0], ";")
      receiveU := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-param-recipient-user-name=")[1], ";")[0], ";")
      receiveD := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-param-recipient-display-name=")[1], ";")[0], ";")
      months := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-param-months=")[1], ";")[0], ";")
      plan := strings.TrimRight(strings.SplitAfter(strings.SplitAfter(m, "msg-param-sub-plan=")[1], ";")[0], ";")
      b.Subgift(giftU, giftD, receiveU, receiveD, months, plan)
    }
  }
}

func (b *Bot) SendMsg(msg string) {
  if Query("SELECT state FROM module where moduleName = 'sendMe'") == "1" {msg = "/me "+msg}
  fmt.Fprintf(b.C, "PRIVMSG #%s :"+msg+"\r\n", b.Channel,)
  fmt.Println("[SEND] " + msg)
}

func (b *Bot) SendWhisper(msg, user string) {
  fmt.Fprintf(b.C, "PRIVMSG #%s :/w "+user+" "+msg+"\r\n", b.Channel,)
  fmt.Println("[WHISPER] "+user+": "+msg)
}

func (b *Bot) SendTimeout(user, time, msg string) {
  fmt.Fprintf(b.C, "PRIVMSG #%s :/timeout "+user+" "+time+" "+msg+"\r\n", b.Channel,)
  fmt.Println("[TIMEOUT] " + msg)
}
