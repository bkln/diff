package main

import "fmt"
import "puzz"
import "time"

import (
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
)

func Server() {
	fmt.Println("Starting server...")
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(puzz.DupOrInsertService), "")
	http.Handle("/rpc", s)
	http.ListenAndServe(":10000", nil)
}

func Ramp() {
	for i := 0; i < 10000; i++ {
		cvec := puzz.GetRandCVec()
		_, _ = puzz.HasPic(cvec)
		err := puzz.AddPic(cvec, fmt.Sprintf("%09d", i))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func Test() {
	images := []string{"lena.jpg", "lena2.jpg", "pjw.jpg"}
	for _, image := range images {
		cvec, err := puzz.GetCVec(image)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		for i := 0; i < 100; i++ {
			_, _ = puzz.HasPic(cvec)
		}
		duplicates, err := puzz.HasPic(cvec)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if len(duplicates) == 0 {
			fmt.Println("No dupliactes for", image)
			err := puzz.AddPic(cvec, image)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Printf("Found duplicates for %s: %s\n", image, duplicates)
		}
		fmt.Println(time.Now())
	}
}

func main() {
	fmt.Println(time.Now())
	err := puzz.InitDb()
	if err != nil {
		fmt.Println("Error initializing db:", err)
		return
	}
	fmt.Println(time.Now())
	Server()

	Ramp()
	Test()
	fmt.Println(time.Now())
}
