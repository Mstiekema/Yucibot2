package base

import (
  
)

func (b *Bot) UpdateLines(U User) {
  StringOldPoints := Query("user", "points", U.username)
  if StringOldPoints != "" {
    Update("user", "num_lines = num_lines + 1", "'"+U.username+"'")
  } else {
    Insert("user (name, userId, points, num_lines, level, isMod)", "('"+ U.username + "', '" + U.userId + "', '0', '1', '100', '"+ U.mod + "')")
  }
}