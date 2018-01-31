package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
)

type Message struct {
	XMLName  xml.Name `json:"-" xml:"msg"`
	Sender   string   `json:"txer" xml:"txer"`
	Receiver string   `json:"rxer" xml:"rxer"`
	Content  string   `json:"msg" xml:"content"`
}

func main() {
	m := Message{
		Sender:   "Alice",
		Receiver: "Bob",
		Content:  "Yo.",
	}

	buf := &bytes.Buffer{}

	//encoder := json.NewEncoder(buf)
	encoder := xml.NewEncoder(buf)
	if err := encoder.Encode(m); err != nil {
		log.Fatal(err)
	}

	if err := encoder.Encode(struct {
		XMLName xml.Name `json:"-" xml:"wrapper"`
		Boo     string   `json:"boo" xml:"ahhh"`
	}{
		Boo: "Ahhhh!!!!",
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("json: %s\n", buf)
}
