package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// We are using generics, K is any type that is comparable so that we can perform equality and relational operations.
type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) (V, error)
}

// KVStore is succesfully implementing the Storer interface because it implements all the methods mentioned in the interface.
type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// *KVStore[K, V] indicates that the function returns a pointer to a Storer instance.
// &KVStore[K, V] line creates a new instance of KVStore and returns its address.
// The & operator is used to get the address of the newly created Storer instance.
// NewKVStore is a Constructor Function, it creates and initializes a new KVStore instance.
func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

// Note: Has function is not concurrent safe, should be used with a lock/mutex.
func (s *KVStore[K, V]) Has(key K) bool {
	_, ok := s.data[key]
	return ok
}

// Put is a method defined on the KVStore struct
func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}

	return value, nil
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Has(key) {
		return fmt.Errorf("the key (%v) does not exist", key)
	}
	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Delete(key K) (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}

	delete(s.data, key)

	return value, nil
}

// type Server struct {
// 	Store Storer[string, string]
// }

// func (s *Server) getUserByName(name string) (string, error) {
// 	return s.Store.Get(name)
// }

// func main() {
// 	s := Server{
// 		Store: NewKVStore[string, string](),
// 	}

// 	if err := s.Store.Put("Fuck U", "BITCHH!!"); err != nil {
// 		log.Fatal(err)
// 	}

// 	store := NewKVStore[string, string]()

// 	if err := store.Put("Aryan", "Kobe"); err != nil {
// 		log.Fatal(err)
// 	}

// 	value, err := store.Get("Aryan")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := store.Update("Aryan", "MAMBA24"); err != nil {
// 		log.Fatal(err)
// 	}
// 	// Below we cannot use `:=`, it will throw error, as `:=` is used to declare and infere the type of a value
// 	// and store it in the variable on the left hand side, where as `=` should be used once we have already declared
// 	// a variable.
// 	value, err = store.Get("Aryan")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(value)

// }

type Server struct {
	Storage    Storer[string, string]
	ListenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		Storage:    NewKVStore[string, string](),
		ListenAddr: listenAddr,
	}
}

// // Basic HTTP server, without using any external frameworks listening on port 3000
// func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("picabooo!"))
// }

// func (s *Server) Start() {
// 	fmt.Printf("HTTP server is running on port %s", s.ListenAddr)

// 	http.HandleFunc("/put", s.handlePut)

// 	log.Fatal(http.ListenAndServe(s.ListenAddr, nil))
// }

// Using the echo web framework.
func (s *Server) handlePut(c echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	s.Storage.Put(key, value)

	return c.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}

func (s *Server) handleGet(c echo.Context) error {
	key := c.Param("key")

	value, err := s.Storage.Get(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"value": value})
}

func (s *Server) handleUpdate(c echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	s.Storage.Update(key, value)

	return c.JSON(http.StatusOK, map[string]string{"updated-value": value})
}

func (s *Server) handleDelete(c echo.Context) error {
	key := c.Param("key")

	s.Storage.Delete(key)

	return c.JSON(http.StatusOK, map[string]string{"deleted-entry": key})
}

func (s *Server) Start() {
	fmt.Printf("HTTP server is running on port %s", s.ListenAddr)

	e := echo.New()

	e.GET("/put/:key/:value", s.handlePut)
	e.GET("/get/:key", s.handleGet)
	e.GET("/update/:key/:value", s.handleUpdate)
	e.GET("/delete/:key", s.handleDelete)

	e.Start(s.ListenAddr)
}

func main() {
	s := NewServer(":3000")
	s.Start()
}
