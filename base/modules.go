package base

import (
  "strings"
)

func (b *Bot) Modules(C string, usr string) {
  b.UserInfoComms(C, usr)
  b.Basic(C, usr)
}

func (b *Bot) UserInfoComms(msg string, usr string) {
  if msg == "!points" {
    var res = Query("user", "points", strings.ToLower(usr))
    if res != "" {b.SendMsg(usr + " currently has " + res + " points")}
  } 
  if msg == "!lines" {
    var res = Query("user", "num_lines", strings.ToLower(usr))
    if res != "" {b.SendMsg(usr + " currently has written " + res + " lines in chat")}
  }
}

func (b *Bot) Basic(msg string, usr string) {
  if msg == "!test" {
    b.SendMsg("This is a test message created by " + usr + " FeelsGoodMan")
  } 
}