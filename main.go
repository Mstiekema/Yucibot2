package main

import (
  "github.com/Mstiekema/Yucibot2/base"
  "github.com/Mstiekema/Yucibot2/web"
)

func main() {  
  bot := base.CrtBot()    
  bot.Connect()
  
  go web.MainWeb()
  bot.Reader(bot.C)
}