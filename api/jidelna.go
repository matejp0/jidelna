package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const URL = "https://www.jidelna.cz/rest/u/c58zbtfnjz72h6t5nzfva9uzvbag8m/"

type User struct {
  userId string
  cookies []*http.Cookie
}


func (u *User) GetUserInfo() string {
  client := http.Client{}

  req, _ := http.NewRequest(http.MethodGet, URL + "/uzivatel/" + u.userId + "/info", nil)
  req.AddCookie(u.cookies[1])
  response, err := client.Do(req)
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
  u.cookies = httpResponse.Cookies()
  u.userId = n
}

func GetFoods() {
  t := time.Now()
  httpResponse, err := http.Get(URL+"zarizeni/356/dny/od/" + t.Format("2006-01-02") + "/do/" + t.AddDate(0, 0, 1).Format("2006-01-02"))
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  body, err := ioutil.ReadAll(httpResponse.Body)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
  
  var parsedResponse []Day

  err = json.Unmarshal(body, &parsedResponse)
  if err != nil {
    log.Fatal(err)
  }

  PrintFoods(parsedResponse)
}

func PrintFoods(days []Day) {
  for _, day := range days {
    fmt.Println(day.Date)
    for _, item := range day.Den.CastiDne[0].Menu {
      if item.LzeObjednat == false { continue }
      fmt.Printf("[%v]\n", item.Nazev)
      for _, values := range item.Chody {
        switch values.Nazev {
        case "Polévka":
          fmt.Printf("\t%v: %v\n", values.Nazev, values.Jidlo)
        case "Jídlo":
          fmt.Printf("\t%v: %v\n", values.Nazev, values.Jidlo)
        case "Příloha":
          fmt.Printf("\t%v: %v\n", values.Nazev, values.Jidlo)

        }
      }
    }
    fmt.Println()
  }
}

type Day struct {
  Date string `json:"datum"`
  Den struct { 
    CastiDne []struct {
      Nazev string `json:"nazev"` // "oběd" -- velice hodnotná informace
      Menu []struct {
        Nazev string `json:"nazev"`
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
    Ucty map[string]interface{} `json:"ucty"`
  }
}
