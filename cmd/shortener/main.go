package main

import (
	"go-shortener/handlers"
	"log"
	"net/http"
)

//var DB = app.InitDB()

func main() {

	//DB.Insert("raw-1", "short-1")
	//fmt.Println(DB)
	//post, err := DB.Get("short-short-short")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(post.RawUrl)
	//}
	mux := http.NewServeMux()
	//mux.Handle("/", &handlers.MyHandler{})
	mux.HandleFunc("/", handlers.MyHandler)
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
