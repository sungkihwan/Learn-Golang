package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("잔액이 부족합니다")

// 계좌 생성
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account // 새로운 객체를 같은 메모리를 참조하여 반환(포인터)
}

// 입금
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// 잔액
func (a Account) Balance() int {
	return a.balance
}

// 인출
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil // null
}

func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

func (a Account) Owner() string {
	return a.owner
}

// 내장 tostring 변환
func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance())
}
