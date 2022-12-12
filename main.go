package main

import (
	"fmt"
	"time"
)

type user struct {
	id    int
	name  string
	timer *time.Timer
	stop  chan struct{}
}
type users map[int]*user

var us users

const userNum = 10
const userPrefix = "user"

func init() {
	us = make(users, userNum)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
	fmt.Scanln()
}

func run() error {
	for i := 0; i < userNum; i++ {
		newUser(i + 1)
		// wait
		time.Sleep(1 * time.Second)
	}

	time.Sleep(userNum / 2 * time.Second)

	//　stop remaining timers
	for _, u := range us {
		us.stop(u.id)
	}
	return nil
}

func newUser(id int) {
	name := userPrefix + fmt.Sprintf("%02d", id)
	t := time.NewTimer(userNum * time.Second)
	u := &user{
		id:    id,
		name:  name,
		timer: t,
		stop:  make(chan struct{}),
	}
	// regist new user
	us[id] = u
	fmt.Printf("[add]id: %d\n", id)

	// timer's wait
	go u.wait()
}

func (u *user) wait() {
	select {
	case <-u.timer.C:
	case <-u.stop:
		fmt.Println("received stop signal")
	}

	u.clean()
	fmt.Println("goroutine complete")
}

func (u *user) clean() {
	fmt.Printf("[delete]id: %d, name: %s\n", u.id, u.name)
	delete(us, u.id)
}

func (u users) stop(id int) {
	user, ok := u[id]
	if !ok {
		fmt.Println("")
	}

	if user.timer.Stop() {
		fmt.Println("forced stop")
		// trigger on cleaning
		user.stop <- struct{}{}
	} else {
		<-user.timer.C
		fmt.Println("already stopped")
	}
}
