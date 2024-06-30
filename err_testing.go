package main

import "log"

type myErr struct {
	msg string
}

func (e *myErr) Error() string {
	return e.msg
}

func getInfo(isCat bool) error {
	var err *myErr
	var e error

	e = err
	if e != nil {
		log.Println("err != nil")
	} else {
		log.Println("err == nil")
	}

	if isCat {
		return err
	} else {
		return nil
	}
}

func ErrText() {
	err := getInfo(true)
	if err != nil {
		log.Println("err != nil")
	} else {
		log.Println("err == nil")
	}

	err = getInfo(false)
	if err != nil {
		log.Println("err != nil")
	} else {
		log.Println("err == nil")
	}
}
