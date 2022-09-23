package main

import (
	"errors"
	"fmt"
	"sync"
)

type ChatRoom struct {
	sync.RWMutex
	userPool map[string]*User
}

var (
	singleChatRoom *ChatRoom
	createRoomOnce sync.Once
)

func getChatRoom() *ChatRoom {
	createRoomOnce.Do(func() {
		singleChatRoom = &ChatRoom{
			userPool: make(map[string]*User),
		}
	})

	return singleChatRoom
}

func (r *ChatRoom) addUser(user *User) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.userPool[user.name]; ok {
		return errors.New("the username is already taken")
	}

	r.userPool[user.name] = user

	return r.broadcast(
		user, []byte(fmt.Sprintf("* %s has entered the chat\n", user.name)),
	)
}

func (r *ChatRoom) broadcast(sender *User, msg []byte) error {
	fmt.Print(string(msg))
	for _, user := range r.userPool {
		if user.name == sender.name {
			continue
		}

		if _, err := user.conn.Write(msg); err != nil {
			return err
		}
	}

	return nil
}

func (r *ChatRoom) getUserNames() []string {
	r.RLock()
	defer r.RUnlock()

	var names []string
	for name := range r.userPool {
		names = append(names, name)
	}

	return names
}

func (r *ChatRoom) removeUser(user *User) {
	r.Lock()
	defer r.Unlock()

	delete(r.userPool, user.name)
	r.broadcast(user, []byte(fmt.Sprintf("* %s has left the chat\n", user.name)))
}
