package crudgolab // import "gopkg.in/kravchadev/crudgolab.v1"

import (
	"encoding/json"
	"net/url"
)

type StorageStatusResponse interface {
	Error() string   // ошибки
	StatusCode() int // коды статусов http
}

type Storage interface {
	// создает ресурс и сохраняет в нем данные, затем возвращает идентификатор
	Create(collection string, body *json.Decoder, query url.Values) (string, StorageStatusResponse)

	// читаем ресурс по ид
	Get(collection, id string, query url.Values) (interface{}, StorageStatusResponse)

	// читае все ресурсы в указанной коллекции
	GetAll(collection string, query url.Values) ([]interface{}, StorageStatusResponse)

	// обновить ресурс
	Update(collection, id string, body *json.Decoder, query url.Values) StorageStatusResponse

	// удалить ресурс
	Delete(collection, id string, query url.Values) StorageStatusResponse

	// удалить все ресурсы в коллекции
	DeleteAll(collection string, query url.Values) StorageStatusResponse
}
