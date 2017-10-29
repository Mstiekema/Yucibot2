package base

import (
  
)

func (b *Bot) Modules(C string, U User) {
  b.UserInfoComms(C, U)
  b.Basic(C, U)
  b.Roulette(C, U)
}

func (b *Bot) UserInfoComms(C string, U User) {
  if C == "!points" {
    var res = Query("user", "points", U.username)
    if res != "" {b.SendMsg(U.displayName + " currently has " + res + " points")}
  } 
  if C == "!lines" {
    var res = Query("user", "num_lines", U.username)
    if res != "" {b.SendMsg(U.displayName + " currently has written " + res + " lines in chat")}
  }
}

func (b *Bot) Basic(C string, U User) {
  if C == "!test" {
    b.SendMsg("This is a test message created by " + U.displayName + " FeelsGoodMan")
  }
}