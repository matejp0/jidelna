package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/matejp0/jidelna/api"

	"github.com/peterbourgon/ff/v3/ffcli"
)

const DATE_J = "2006-01-02"
const DATE_H = "02. 01. 2006"

func main() {
  var user api.User

  user.Login(EMAIL, PASSWORD)

  var (
    listFlagSet = flag.NewFlagSet("jidelna list", flag.ExitOnError)
    nDays = listFlagSet.Int("n", 7, "How many days to list")
    orderFlagSet = flag.NewFlagSet("jidelna order", flag.ExitOnError)
    date = orderFlagSet.String("d", "", "Date")
  )

  list := &ffcli.Command{
    Name: "list",
    ShortUsage: "jidelna list [-n days]",
    ShortHelp: "Lists foods",
    FlagSet: listFlagSet,
    Exec: func(ctx context.Context, args []string) error {
      t := time.Now()
      PrintFoods(user.GetFoods(t.Format(DATE_J), t.AddDate(0, 0, *nDays).Format(DATE_J)), user)
      return nil
    },
  }

  info := &ffcli.Command{
    Name: "info",
    ShortUsage: "jidelna info",
    ShortHelp: "Get user info",
    Exec: func(ctx context.Context, args []string) error {
      fmt.Println(user.GetUserInfo())
      return nil
    },
  }

  order := &ffcli.Command{
    Name: "order",
    ShortUsage: "jidelna order [-d date] <id>",
    ShortHelp: "Order food",
    FlagSet: orderFlagSet,
    Exec: func(ctx context.Context, args []string) error {
      if len(args) < 1 {
        return fmt.Errorf("Requires exactly 1 argument")
      }
      num, err := strconv.Atoi(args[0])
      if err != nil {
        return fmt.Errorf("Failed to convert food id to int")
      }

      if len(*date) == 0 {
        return fmt.Errorf("You must specify the date")
      } else {
        fmt.Println(len(*date), *date)
      }

      success := user.EditFood(num, *date)
      PrintFoods(user.GetFoods(*date, *date), user)
      if success {
        fmt.Println("Successfully ordered")
        return nil
      }
      return fmt.Errorf("Failed to edit the food")
      
    },
  }
  root := &ffcli.Command {
    ShortUsage: "jidelna <subcommand>",
    Subcommands: []*ffcli.Command{list, info, order},
  }

  if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err.Error())
  }
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

