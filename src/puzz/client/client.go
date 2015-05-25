package main

import "puzz"
import "net/rpc/jsonrpc"
import "net/rpc"
import "fmt"
import "net"
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	c *rpc.Client
)

type Req struct {
	Method string                `json:"method"`
	Id     int                   `json:"id"`
	Params []puzz.DupOrInsertArg `json:"params"`
}

func Req2() error {
	r := &Req{
		Method: "DupOrInsertService.DupOrInsert",
		Id:     1,
		Params: []puzz.DupOrInsertArg{
			puzz.DupOrInsertArg{
				Name: "Testimage",
				Cvec: puzz.GetRandCVec(),
			},
		},
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("Marshal: %v", err)
		return err
	}
	//fmt.Println("req:", string(data))
	resp, err := http.Post("http://127.0.0.1:10000/rpc",
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Printf("Post: %v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %v", err)
		return err
	}
	//fmt.Println("resp: ", string(body))
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
		return err
	}
	//fmt.Println(result)
	return nil

}

func main() {
	conn, err := net.Dial("tcp", "localhost:10000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	c = jsonrpc.NewClient(conn)

	count := 0

	ci := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				Req2()
				if err != nil {
					fmt.Println(err)
				}
				ci <- 1
			}
		}()
	}
	fmt.Println("Count:", count, "; Time:", time.Now())
	for i := range ci {
		count += i
		if count%1000 == 0 {
			fmt.Println("Count:", count, "; Time:", time.Now())
		}
	}

	return
	for {
		err = Req2()
		if err != nil {
			fmt.Println(err)
		}
		count++
		if count%1000 == 0 {
			fmt.Println("Count:", count, "; Time:", time.Now())
		}
	}
}
