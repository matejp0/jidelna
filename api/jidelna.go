package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const URL = "https://www.jidelna.cz/rest/u/c58zbtfnjz72h6t5nzfva9uzvbag8m/"

type User struct {
  UserId string
  schoolId string
  client http.Client
}


func (u *User) GetUserInfo() UserInfo {
  req, _ := http.NewRequest(http.MethodGet, URL + "/uzivatel/" + u.UserId + "/info", nil)
  response, err := u.client.Do(req)
  if err != nil {
    log.Fatal(err)
  }
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Fatal(err)
    return UserInfo{}
  }

  var userInfo UserInfo
  err = json.Unmarshal(body, &userInfo)
  
  if err != nil {
    log.Fatal(err)
    return UserInfo{}
  }

  

  return userInfo
}

func (u *User) Login(email, password string) bool {
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
    return false
  }

  var user LogInUser

  err = json.Unmarshal(body, &user)
  if err != nil {
    log.Fatal(err)
    return false
  }

  // this is stupid but it's because of the stupidity of jidelna
  var n string
  for i := range user.Ucet.Ucty { 
    n = i
  }
  urlObj, _ := url.Parse(URL)
  u.client.Jar.SetCookies(urlObj, httpResponse.Cookies()[1:]) 
  u.UserId = n
  u.schoolId = user.Ucet.Ucty[n]["regc"].(string)

  return true
}

func (u *User) EditFood(idMenu int, date string) bool {
  food := u.createFood(idMenu, date)
  marshalized, err := json.Marshal([1]Food{food})
  if err != nil {
    log.Fatal("Marshalization failed", err)
  }
  jsonList := make([]string, 0)
  jsonList = append(jsonList, string(marshalized[:]))
  urlValues := url.Values{}
  urlValues.Set("json", string(marshalized[:]))
  req, _ := http.NewRequest(http.MethodPost, URL + "zarizeni/" + u.schoolId + "/objednavky", strings.NewReader(urlValues.Encode()))
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  resp, err := u.client.Do(req)
  if err != nil {
    log.Fatal(err)
    return false
  }
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

  var jidelnaResponse JidelnaResponse
  err = json.Unmarshal(body, &jidelnaResponse)

  if err != nil {
    log.Fatal(err)
    return false
  }

  if jidelnaResponse.Stav == "ok" { 
    return true
  } else {
    return false
  }

}

func (u *User) createFood(idMenu int, date string) Food {
  return Food{
    IdUzivatele: u.UserId,
    IdMenu: strconv.Itoa(idMenu),
    Den: date,
    Stav: "Prihlaseno",
    Mnozstvi: 1,
  }
}

func (u *User) GetFoods(startDate string, endDate string) []Day {
  req, _ := http.NewRequest(http.MethodGet, URL+"zarizeni/" + u.schoolId + "/dny/od/" + startDate + "/do/" + endDate, nil)

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

  return parsedResponse
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
    Ucty map[string]map[string]any `json:"ucty"`
  }
}


type Food struct {
  IdUzivatele string `json:"idUzivatele"`
  IdMenu string `json:"idMenu"`
  Den string `json:"den"`
  Stav string `json:"stav"`
  Mnozstvi int `json:"mnozstvi"`
}

type JidelnaResponse struct {
  Stav string `json:"stav"`
}

type UserInfo struct {
  Jmeno string `json:"jmeno"`
  KontoProObjednavani string `json:"kontoProObjednavani"`
}
