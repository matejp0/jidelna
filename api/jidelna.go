package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

const URL = "https://www.jidelna.cz/rest/u/c58zbtfnjz72h6t5nzfva9uzvbag8m/"

type User struct {
  userId string
  client http.Client
}


func (u *User) GetUserInfo() string {
  req, _ := http.NewRequest(http.MethodGet, URL + "/uzivatel/" + u.userId + "/info", nil)
  response, err := u.client.Do(req)
  if err != nil {
    log.Fatal(err)
  }
  responseDecoded, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Fatal(err)
  }
  return string(responseDecoded)
}

func (u *User) Login(email, password string) () {
  jar, err := cookiejar.New(nil)
  if err != nil {
    log.Fatal(err.Error())
  }
  u.client = http.Client{
    Jar: jar,
  }

  data := url.Values{}
  data.Add("login", email)
  data.Add("heslo", password)
  httpResponse, err := http.PostForm(URL+"login/jmenoheslo", data)
  if err != nil {
    log.Fatal("Failed to log in", err)
    os.Exit(1)
  }
  body, err := ioutil.ReadAll(httpResponse.Body)
  if err != nil {
    log.Fatal(err)
  }

  var user LogInUser

  err = json.Unmarshal(body, &user)
  if err != nil {
    log.Fatal(err)
  }

  // this is stupid but it's because of the stupidity of jidelna
  var n string
  for i := range user.Ucet.Ucty { 
    n = i
  }
  urlObj, _ := url.Parse(URL)
  u.client.Jar.SetCookies(urlObj, httpResponse.Cookies()[1:]) 
  u.userId = n
}

func (u *User) GetFoods(days int) {
  t := time.Now()
  req, _ := http.NewRequest(http.MethodGet, URL+"zarizeni/356/dny/od/" + t.Format("2006-01-02") + "/do/" + t.AddDate(0, 0, days).Format("2006-01-02"), nil)

  resp, err := u.client.Do(req)

  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  var parsedResponse []Day

  err = json.Unmarshal(body, &parsedResponse)
  if err != nil {
    log.Fatal(err)
  }

  PrintFoods(parsedResponse, u)
}

func PrintFoods(days []Day, u *User) {
  d := color.New(color.FgCyan, color.Bold)
  dateColor := color.New(color.FgMagenta, color.Italic)
  for _, day := range days {
    date, _ := time.Parse("2006-01-02", day.Date)

    dateColor.Println(date.Format("01. 02. 2006"))
    
    ucet := day.Den.CastiDne[0].Objednavky[u.userId].(map[string]any)
    for _, item := range day.Den.CastiDne[0].Menu {
      if item.LzeObjednat == false { continue }
      if strconv.Itoa(item.Id) == ucet["idMenu"] && ucet["stav"] == "Prihlaseno" {
        d = color.New(color.FgHiYellow, color.Bold)
      } else if ucet["stav"] == "Vyzvednuto"{
        d = color.New(color.FgHiRed)
      } else {
        d = color.New(color.FgHiWhite)
      }

      d.Printf("[%v]", item.Nazev)
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

type Day struct {
  Date string `json:"datum"`
  Den struct { 
    CastiDne []struct {
      Objednavky map[string]any `json:"objednavky"`
      Nazev string `json:"nazev"` // "oběd" -- velice hodnotná informace
      Menu []struct {
        Nazev string `json:"nazev"`
        Id int `json:"id"`
        LzeObjednat bool `json:"lzeObjednat"`
        Chody []struct {
          Nazev string `json:"nazev"`
          Jidlo string `json:"jidlo"`
        }
      }
    }
  }
}

type LogInUser struct {
  Ucet struct {
    Ucty map[string]any `json:"ucty"`
  }
}
