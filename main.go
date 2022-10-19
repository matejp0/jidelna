package main

import (
	"github.com/matejp0/jidelna/api"
)

func main() {
  var user api.User
  user.Login(EMAIL, PASSWORD)
  user.GetFoods()

    
 // fmt.Println(user.GetUserInfo())
}
