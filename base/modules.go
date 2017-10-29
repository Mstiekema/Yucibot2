package base

import (
  "strings"
)

func (b *Bot) Modules(C, usr, M string) {
  b.UserInfoComms(C, usr)
  b.Basic(C, usr)
  b.Roulette(C, usr, M)
}

func (b *Bot) UserInfoComms(C, usr string) {
  if C == "!points" {
    var res = Query("user", "points", strings.ToLower(usr))
    if res != "" {b.SendMsg(usr + " currently has " + res + " points")}
  } 
  if C == "!lines" {
    var res = Query("user", "num_lines", strings.ToLower(usr))
    if res != "" {b.SendMsg(usr + " currently has written " + res + " lines in chat")}
  }
}

func (b *Bot) Basic(C, usr string) {
  if C == "!test" {
    b.SendMsg("This is a test message created by " + usr + " FeelsGoodMan")
  } 
}