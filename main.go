package main

import (
	"fmt"
	"time"
)

type user struct {
	id    int
	name  string
	timer *time.Timer
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
		go add(i)
		// wait
		time.Sleep(1 * time.Second)
	}
	return nil
}

func add(n int) {
	t := time.NewTimer(5 * time.Second)
	users[n] = &user{
		id:    n,
		name:  names[n],
		timer: t,
	}
	fmt.Printf("[add]id: %d\n", n)

	// timer wait
	<-users[n].timer.C

	fmt.Println("stop")
	clean(n)
}

func clean(n int) error {
	fmt.Printf("[delete]id: %d, name: %s\n", n, users[n].name)
	delete(users, n)
	return nil
}
