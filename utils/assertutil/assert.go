package assertutil

import (
	"fmt"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Assert(isOk bool, err error) {
	if !isOk {
		panic(err)
	}
}

func AssertStringError(isOk bool, str string) {
	if !isOk {
		panic(fmt.Errorf("%s", str))
	}
}

func WithRecover(work func() error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic:", err)
		}
	}()

	err := work()
	if err != nil {
		fmt.Println(err)
	}
}
