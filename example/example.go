package main

import (
	"log"
	"net/http"

	"github.com/kravchadev/crudgolab"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("хэллоу ворлд!\n"))
}

func main() {
	// инициализация коллекций
	storage := NewMapStorage()
	storage.AddMap("artists")
	storage.AddMap("albums")

	// создаем роуты
	api := crudgolab.New(storage)

	// создаем endpoint для api
	http.Handle("/api/", http.StripPrefix("/api", api))

	// монтируем обработчик по дефолту (если нет никакой записи после api/)
	http.HandleFunc("/", hello)

	// старт сервера
	log.Println("сервер слушает запросы на localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
