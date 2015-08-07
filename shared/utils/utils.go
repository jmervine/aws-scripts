package utils

import (
	"fmt"
	"os"
	"os/user"
	"sync"
	"time"
)

const DAY = time.Second * 60 * 60 * 24

func Home() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func DaysAgo(i int, from *time.Time) *time.Time {
	if from == nil {
		now := time.Now()
		from = &now
	}

	if i == 0 {
		return from
	}

	ret := (*from).Add(time.Duration(i*-1) * DAY)
	return &ret
}

var locker sync.Mutex

type Output struct {
	File *os.File
	Lock bool // support simple locking for thread safty
}

func NewOutput(p string, l bool) *Output {
	o := &Output{Lock: l}

	f, e := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil {
		panic(e)
	}
	o.File = f

	return o
}

// to create output with a lock or file, simply do:
//
// output := new(Output)

func (o *Output) Puts(s string) {
	if o.Lock {
		locker.Lock()
	}

	if o.File != nil {
		if _, e := o.File.WriteString(s); e != nil {
			panic(e)
		}
		if e := o.File.Sync(); e != nil {
			panic(e)
		}
	}

	fmt.Printf(s)

	if o.Lock {
		locker.Unlock()
	}
}
