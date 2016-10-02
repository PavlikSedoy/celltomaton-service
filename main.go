package main

import (
	"encoding/json"
	"fmt"
	"github.com/initiumsrc/celltomaton"
	"log"
	"net/http"
	"os"
)

func main() {
	f, err := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("testing...")

	http.HandleFunc("/", cellHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

type CellConfig struct {
	Array  []int
	Height int
	Rule   int
}

func (c *CellConfig) containsNil() bool {
	if c.Array == nil || c.Height == -1 || c.Rule == -1 {
		return true
	}
	return false
}

type MsgResp struct {
	Matrix [][]int
}

func cellHandler(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			resp.Header().Set("status", "500")
			fmt.Fprintln(resp, "Panic: ", r)
		}
		log.Println(req.Method + ": " + req.URL.String())
	}()

	resp.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method != "POST" {
		resp.Header().Set("status", "405")
		fmt.Fprintln(resp, "405 only takes post")
		return
	}

	cellc := &CellConfig{Array: nil, Height: -1, Rule: -1}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&cellc)

	if err != nil || cellc.containsNil() {
		resp.Header().Set("status", "400")
		fmt.Fprintln(resp, "400 must supply conformant json")
		fmt.Fprintln(resp, "\nExample:\n{\n \"array\": [0, 1, 1],\n \"height\": 10,\n \"rule\": 3\n}")
		return
	}

	matrix := celltomaton.Get(cellc.Array, cellc.Height, cellc.Rule)
	b, err := json.Marshal(&MsgResp{Matrix: matrix})

	if err != nil {
		panic("Matrix was not json conformant, internal error")
	}

	resp.Write(b)
}
