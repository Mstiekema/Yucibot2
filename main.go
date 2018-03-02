package main

import (
  "time"
  "github.com/Mstiekema/Yucibot2/base"
  "github.com/Mstiekema/Yucibot2/web"
)

func main() {  
  bot := base.CrtBot()
  bot.Connect()
  
  go web.MainWeb()
  go bot.Reader(bot.C)
  
  t := time.NewTicker(300 * time.Second)
  for {
    bot.UpdateUser()
    <-t.C
  }
}