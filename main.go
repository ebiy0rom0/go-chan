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

// user list
var users map[int]*user

var names = []string{
	"hoge",
	"piyo",
	"fuga",
	"var",
	"lib",
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
	fmt.Scanln()

	fmt.Printf("%+#v", users)
}

func run() error {
	users = make(map[int]*user, 5)
	for i := 0; i < 5; i++ {
		add(i)
		go wait(i)

		// wait
		time.Sleep(1 * time.Second)
	}

	time.Sleep(2 * time.Second)

	for _, u := range users {
		if u.timer.Stop() {
			fmt.Println("forced stop")
			// trigger to cleaning
			u.stop <- struct{}{}
		} else {
			<-u.timer.C
			fmt.Println("already stopped")
		}
	}
	return nil
}

func add(n int) {
	t := time.NewTimer(5 * time.Second)
	users[n] = &user{
		id:    n,
		name:  names[n],
		timer: t,
		stop:  make(chan struct{}),
	}
	fmt.Printf("[add]id: %d\n", n)
}

func wait(n int) {
	select {
	case <-users[n].timer.C:
	case <-users[n].stop:
		fmt.Println("received stop signal")
	}

	clean(users[n].id)
	fmt.Println("fin")
}

func clean(n int) error {
	fmt.Printf("[delete]id: %d, name: %s\n", n, users[n].name)
	delete(users, n)
	return nil
}
