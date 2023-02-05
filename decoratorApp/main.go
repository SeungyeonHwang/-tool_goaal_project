package main

import (
	"fmt"

	"github.com/SeungyeonHwang/tool-goaal/decoratorApp/cipher"
	"github.com/SeungyeonHwang/tool-goaal/decoratorApp/lzw"
)

type Component interface {
	Operator(string)
}

var sentData string
var receiveData string

type SendComponent struct{}

func (self *SendComponent) Operator(data string) {
	//send data
	sentData = data
}

// decorator component(Zip)
type ZipComponent struct {
	key string
	com Component
}

// Zip
func (self *ZipComponent) Operator(data string) {
	zipData, err := lzw.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(zipData))
}

// decorator component(Encrypt)
type EncryptComponent struct {
	key string
	com Component
}

// Encrypt
func (self *EncryptComponent) Operator(data string) {
	encrptData, err := cipher.Encrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(encrptData))
}

// Decrpt
type DecryptComponent struct {
	key string
	com Component
}

func (self *DecryptComponent) Operator(data string) {
	decrpytData, err := cipher.Decrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(decrpytData))
}

// decorator component(unZip)
type unZipComponent struct {
	key string
	com Component
}

// Zip
func (self *unZipComponent) Operator(data string) {
	unZipData, err := lzw.Read([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unZipData))
}

type ReadComponent struct{}

func (self *ReadComponent) Operator(data string) {
	//send data
	receiveData = data
}

func main() {
	sender := &EncryptComponent{
		key: "abcde",
		com: &ZipComponent{
			com: &SendComponent{},
		},
	}
	sender.Operator("Hello World")

	fmt.Println(sentData)

	receiver := &unZipComponent{
		com: &DecryptComponent{
			key: "abcde",
			com: &ReadComponent{},
		},
	}
	receiver.Operator(sentData)
	fmt.Println(receiveData)
}
