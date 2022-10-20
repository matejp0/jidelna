package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/matejp0/jidelna/api"
)

const DATE_J = "2006-01-02"
const DATE_H = "02. 01. 2006"

func main() {
  var user api.User
  fmt.Println(user.Login(EMAIL, PASSWORD))
  //fmt.Println(user.EditFood(10815, "2022-10-24"))
  t := time.Now()
  PrintFoods(user.GetFoods(t.Format(DATE_J), t.AddDate(0, 0, 5).Format(DATE_J)), user)

    
  fmt.Println(user.GetUserInfo())
}



func PrintFoods(days []api.Day, u api.User) {
  d := color.New(color.FgCyan, color.Bold)
  dateColor := color.New(color.FgMagenta, color.Italic)
  for _, day := range days {
    date, _ := time.Parse(DATE_J, day.Date)

    dateColor.Println(date.Format(DATE_H))
    
    ucet := day.Den.CastiDne[0].Objednavky[u.UserId].(map[string]any)
    for _, item := range day.Den.CastiDne[0].Menu {
      if item.LzeObjednat == false { continue }
      if strconv.Itoa(item.Id) == ucet["idMenu"] && ucet["stav"] == "Prihlaseno" {
        d = color.New(color.FgHiYellow, color.Bold)
      } else if ucet["stav"] == "Vyzvednuto"{
        d = color.New(color.FgHiRed)
        // maybe skip this day entirely?
      } else {
        d = color.New(color.FgHiWhite)
      }

      d.Printf("[%v]", item.Nazev)
      d.Printf(" %v\n", item.Id)
      for _, values := range item.Chody {
        switch values.Nazev {
        case "Polévka":
          d.Printf("\t%v: %v\n", values.Nazev, values.Jidlo)
        case "Jídlo":
          d.Printf("\t%v: %v\n", values.Nazev, values.Jidlo)
        case "Příloha":
          d.Printf("\t%v: %v\n", values.Nazev, values.Jidlo)

        }
      }
      fmt.Println()
    }
    fmt.Println()
  }
}

