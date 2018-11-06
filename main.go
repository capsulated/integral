package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Storage struct {
	Smap 	*sync.Map
}

func (s *Storage) Init() {
	s.Smap = new(sync.Map)
}

func (s *Storage) Set(key string, value interface{}) error {
	if _, ok := s.Smap.Load(key); ok {
		return errors.New("key already exist")
	}
	s.Smap.Store(key, value)
	fmt.Println("value setted")

	go s.delayedCleaning(key)

	return nil
}

func (s *Storage) delayedCleaning(key string) {
	// В данном контексте не вижу смысла в контексте :))))
	// На мой взгляд излишне, здесь вполне подойдёт таймаут
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	select {
	case <-ctx.Done():
		s.Smap.Delete(key)
		fmt.Println("value deleted")
	}
}

func (s *Storage) Get(key string) (interface{}, error) {
	if value, ok := s.Smap.Load(key); ok {
		// Отменять контекст не стал, так как ошибки не произойдёт
		// Хотя можно бы было :)
		s.Smap.Delete(key)
		fmt.Println("value getted and deleted")
		return value, nil
	}
	return 0, errors.New("not found")
}

func main() {
	store := new(Storage)
	store.Init()

	for i := 1; i <= 100; i++ {
		go store.Set("key" + strconv.Itoa(i), i)
	}

	time.Sleep(3 * time.Second)

	// Тест на потокобезопастность
	for i := 1; i <= 20; i++ {
		go func(i int) {
			res, err := store.Get("key" + strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("res: %v\n", res)
		}(i)
	}

	res, err := store.Get("key30")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("res: %v\n", res)

	time.Sleep(15 * time.Second)
}