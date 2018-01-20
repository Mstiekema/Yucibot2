package base

import (
  "fmt"
  "strings"
  "math/rand"
  "github.com/gorilla/websocket"
)

func (b *Bot) Clr(C string, U User) {
  conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9090/post/getCLR/", nil)
  if C == "!clr" && 2 < len(strings.SplitAfter(U.message, " ")) {
    clrType := strings.TrimSpace(strings.SplitAfter(U.message, " ")[1])
    name := strings.TrimSpace(strings.SplitAfter(U.message, " ")[2])
    if err != nil {fmt.Println("dial:", err); return}

    if clrType == "message" {
      exec := func() {
        sendM := []byte(`{"type": "message", "message": "`+U.message+`", "user": "`+U.displayName+`"}`)
        conn.WriteMessage(websocket.TextMessage, sendM)
      }
      b.ExecuteCommand(C, "100", "1000", "30", U, exec)
    } else if clrType == "emote" {
      url := Query("SELECT url FROM emotes WHERE name = '"+name+"'")
      if url == "" {
        b.SendWhisper("This isn't an existing emote", U.username)
      } else {
        exec := func() {
          sendM := []byte(`{"type": "emote", "emote": "`+name+`", "url": "`+url+`"}`)
          conn.WriteMessage(websocket.TextMessage, sendM)
          b.SendWhisper("Succesfully send "+name+" to the stream", U.username)
        }
        b.ExecuteCommand(C, "100", "1000", "30", U, exec)
      }
    } else if clrType == "sound" {
      url := Query("SELECT url FROM clr WHERE type = 'sound' AND name = '"+name+"'")
      if url == "" {
        b.SendWhisper("This isn't an existing sound", U.username)
      } else {
        exec := func() {
          sendM := []byte(`{"type": "sound", "sound": "`+name+`", "url": "`+url+`"}`)
          conn.WriteMessage(websocket.TextMessage, sendM)
          b.SendWhisper("Succesfully send '"+name+"' to the stream", U.username)
        }
        b.ExecuteCommand(C, "100", "1000", "30", U, exec)
      }
    } else if clrType == "gif" {
      url := Query("SELECT url FROM clr WHERE type = 'gif' AND name = '"+name+"'")
      if url == "" {
        b.SendWhisper("This isn't an existing GIF", U.username)
      } else {
        exec := func() {
          sendM := []byte(`{"type": "gif", "gif": "`+name+`", "url": "`+url+`"}`)
          conn.WriteMessage(websocket.TextMessage, sendM)
          b.SendWhisper("Succesfully send '"+name+"' to the stream", U.username)
        }
        b.ExecuteCommand(C, "100", "1000", "30", U, exec)
      }
    } else if clrType == "meme" {
      meme := Query("SELECT url FROM clr WHERE name = '"+name+"'")
      if meme == "" {
        b.SendWhisper("This isn't an existing meme", U.username)
      } else {
        exec := func() {
          sendM := []byte(`{"type": "meme", "meme": "`+meme+`"}`)
          conn.WriteMessage(websocket.TextMessage, sendM)
          b.SendWhisper("Succesfully send '"+meme+"' to the stream", U.username)
        }
        b.ExecuteCommand(C, "100", "2000", "30", U, exec)
      }
    } else {
      b.SendWhisper("This isn't an existing CLR command", U.username)
    }
  }
  if U.message == "!clr meme" || U.message == "!sendmeme" {
    exec := func() {
      db := Conn()
      res, err := db.Query(`SELECT url FROM clr where type = "meme"`)
      if err != nil { panic(err.Error()) }
      defer res.Close()
      
      var urls []string
      for res.Next() {
        var url string
        err = res.Scan(&url)
        urls = append(urls, url)
      }
      conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "meme", "meme": "`+urls[rand.Intn(len(urls))]+`"}`))
      b.SendWhisper("Succesfully send a random meme to the stream", U.username)
      db.Close()
    }
    b.ExecuteCommand(C, "100", "2000", "30", U, exec)
  }
}