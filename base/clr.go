package base

import (
  "strings"
  "github.com/gorilla/websocket"
  "fmt"
)

func (b *Bot) Clr(C string, U User) {
  if C == "!clr" && 2 < len(strings.SplitAfter(U.message, " ")) {
    clrType := strings.TrimSpace(strings.SplitAfter(U.message, " ")[1])
    name := strings.TrimSpace(strings.SplitAfter(U.message, " ")[2])
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9090/post/getCLR/", nil)
    if err != nil {fmt.Println("dial:", err); return}

    if clrType == "message" {
      sendM := []byte(`{"type": "message", "message": "`+U.message+`", "user": "`+U.displayName+`"}`)
      conn.WriteMessage(websocket.TextMessage, sendM)
    } else if clrType == "emote" {
      url := Query("SELECT url FROM emotes WHERE name = '"+name+"'")
      if url == "" {
        b.SendMsg(U.displayName+" this isn't an existing emote")
      } else {
        sendM := []byte(`{"type": "emote", "emote": "`+name+`", "url": "`+url+`"}`)
        conn.WriteMessage(websocket.TextMessage, sendM)
      }
    } else if clrType == "sound" {
      url := Query("SELECT url FROM clr WHERE type = 'sound' AND name = '"+name+"'")
      if url == "" {
        b.SendMsg(U.displayName+" this isn't an existing sound")
      } else {
        sendM := []byte(`{"type": "sound", "sound": "`+name+`", "url": "`+url+`"}`)
        conn.WriteMessage(websocket.TextMessage, sendM)
      }
    } else if clrType == "gif" {
      url := Query("SELECT url FROM clr WHERE type = 'gif' AND name = '"+name+"'")
      if url == "" {
        b.SendMsg(U.displayName+" this isn't an existing GIF")
      } else {
        sendM := []byte(`{"type": "gif", "gif": "`+name+`", "url": "`+url+`"}`)
        conn.WriteMessage(websocket.TextMessage, sendM)
      }
    } else if clrType == "meme" {
      meme := Query("SELECT url FROM clr WHERE name = '"+name+"'")
      if meme == "" {
        b.SendMsg(U.displayName+" this isn't an existing meme")
      } else {
        sendM := []byte(`{"type": "meme", "meme": "`+meme+`"}`)
        conn.WriteMessage(websocket.TextMessage, sendM)
      }
    } else {
      b.SendMsg(U.displayName+" this isn't an existing CLR command")
    }
  }
}