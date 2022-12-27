package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/kravchadev/crudgolab"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type statusResponse struct {
	err        string
	statusCode int
}

func (sr *statusResponse) StatusCode() int { return sr.statusCode }

func (sr *statusResponse) Error() string { return sr.err }

func success(statusCode int) *statusResponse {

	return failure("", statusCode)
}

func failure(err string, statusCode int) *statusResponse {
	return &statusResponse{
		err:        err,
		statusCode: statusCode,
	}
}

func malformedJSON(err error) *statusResponse {
	return failure("malformed JSON: "+err.Error(), http.StatusBadRequest)
}

var (
	collectionNotFound = failure("collection not found", http.StatusNotFound)
	resourceNotFound   = failure("resource not found", http.StatusNotFound)
)

type MapStorage struct {
	*sync.RWMutex
	data map[string]map[string]interface{}
}

// создаем коллекцию
func NewMapStorage() *MapStorage {
	return &MapStorage{&sync.RWMutex{}, make(map[string]map[string]interface{})}
}

// добавляем интерфейс
func (ms *MapStorage) AddMap(collection string) {
	ms.Lock()
	defer ms.Unlock()
	ms.data[collection] = make(map[string]interface{})
}

// возвращаем AddMap().
func (ms *MapStorage) DeleteMap(collection string) {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.data, collection)
}

func (ms *MapStorage) collectionExists(collection string) bool {
	ms.RLock()
	defer ms.RUnlock()
	_, ok := ms.data[collection]

	return ok
}

func (ms *MapStorage) resourceExists(collection, id string) (interface{}, bool) {
	if !ms.collectionExists(collection) {
		return nil, false
	}

	ms.RLock()
	defer ms.RUnlock()
	resource, ok := ms.data[collection][id]

	return resource, ok
}

func (ms *MapStorage) Create(collection string, body *json.Decoder, _ url.Values) (string, crudgolab.StorageStatusResponse) {
	// проверка коллекции на существование
	if !ms.collectionExists(collection) {
		return "", collectionNotFound
	}

	// создаем случайный ID
	id := strconv.FormatInt(rand.Int63(), 10)

	// декод JSON
	var resource map[string]interface{}
	err := body.Decode(&resource)
	if err != nil {
		return "", malformedJSON(err)
	}

	// вставляем данные в коллекцию
	ms.Lock()
	defer ms.Unlock()
	ms.data[collection][id] = resource

	return id, success(http.StatusCreated)
}

func (ms *MapStorage) Get(collection, id string, _ url.Values) (interface{}, crudgolab.StorageStatusResponse) {
	// проверка на существование записи
	resource, ok := ms.resourceExists(collection, id)
	if !ok {
		return nil, resourceNotFound
	}

	return resource, success(http.StatusOK)
}

func (ms *MapStorage) GetAll(collection string, _ url.Values) ([]interface{}, crudgolab.StorageStatusResponse) {
	// проверка на существование коллекции
	if !ms.collectionExists(collection) {
		return nil, collectionNotFound
	}

	// создание среза из всех значений в коллекциях
	var resources []interface{}
	ms.RLock()
	defer ms.RUnlock()
	for _, resource := range ms.data[collection] {
		resources = append(resources, resource)
	}

	return resources, success(http.StatusOK)
}

func (ms *MapStorage) Update(collection, id string, body *json.Decoder, _ url.Values) crudgolab.StorageStatusResponse {
	// проверка на существование записи
	if _, ok := ms.resourceExists(collection, id); !ok {
		return resourceNotFound
	}

	// декод JSON
	var resource map[string]interface{}
	err := body.Decode(&resource)
	if err != nil {
		return malformedJSON(err)
	}

	// обновляем запись
	ms.Lock()
	ms.data[collection][id] = resource
	ms.Unlock()

	return success(http.StatusOK)
}

func (ms *MapStorage) Delete(collection, id string, _ url.Values) crudgolab.StorageStatusResponse {
	// проверка на существование записи
	if _, ok := ms.resourceExists(collection, id); !ok {
		return resourceNotFound
	}

	// удаляем запись
	ms.Lock()
	defer ms.Unlock()
	delete(ms.data[collection], id)

	return success(http.StatusOK)
}

func (ms *MapStorage) DeleteAll(collection string, _ url.Values) crudgolab.StorageStatusResponse {
	// проверка на существование коллекции
	if !ms.collectionExists(collection) {
		return collectionNotFound
	}

	// удаляем коллекцию
	ms.Lock()
	defer ms.Unlock()
	for id := range ms.data[collection] {
		delete(ms.data[collection], id)
	}

	return success(http.StatusOK)
}
