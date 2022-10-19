package main

import (
	"fmt"

	"github.com/matejp0/jidelna/api"
)

func main() {
  api.GetFoods()

  var user api.User
  user.Login(EMAIL, PASSWORD)
  
  fmt.Println(user.GetUserInfo())
}
