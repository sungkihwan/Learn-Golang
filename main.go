package main

import (
	"fmt"

	"github.com/sungkihwan/learngo/accounts"
	"github.com/sungkihwan/learngo/mydict"
)

func main() {
	account := accounts.NewAccount("kihwan")
	account.Deposit(10)
	err := account.Withdraw(15)
	if err != nil {
		// log.Fatalln(err)
		fmt.Println(err)
	}
	fmt.Println(account)

	dictionary := mydict.Dictionary{"name": "kihwan"}
	definition, err := dictionary.Search("see")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
	fmt.Println(dictionary)
}
