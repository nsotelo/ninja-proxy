package main

import "log"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func checkKey(key []byte) {
	for _, size := range []int{16, 24, 32} {
		if len(key) == size {
			return
		}
	}
	log.Fatal("The encryption key has to be 16, 24 or 32 characters in length.")
}

func trapError() {
	r := recover()
	if r != nil {
		log.Printf("Warning: %s", r)
	}
}
