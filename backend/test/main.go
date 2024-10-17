package main

import (
	controllers "drm/handlers"
	"fmt"
)

func main() {
	fmt.Println("go through")
	id := "0e233393-e9f7-4445-a294-80409cb90f76"
	key, _ := controllers.GenerateKey()
	e := controllers.UpdateContentInLCP(id, key, "../uploads/Sway.epub")
	fmt.Println(e)

}
