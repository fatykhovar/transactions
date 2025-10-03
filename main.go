package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

type Transaction struct{
	FromID string
	ToID string
	Amount float64
}

func (u *User) Deposit(sum float64) {
	u.Balance += sum
}

func (u *User) Withdraw(sum float64) error {
	if u.Balance < sum {
		return errors.New("insufficient funds")
	}
	u.Balance -= sum
	return nil
}

func main() {
	user1 := User{ID: "123", Name: "John", Balance: 368.9}
	user2 := User{ID: "456", Name: "Linda", Balance: 698}

	user1.Deposit(80)
	fmt.Println("user1 after deposit 80:", user1)

	if err := user2.Withdraw(80); err != nil {
		fmt.Println("withdraw error:", err)
	}
	fmt.Println("user2 after withdraw 80:", user2)

	if err := user1.Withdraw(500); err != nil {
		fmt.Println("withdraw error:", err)
	}
	fmt.Println("user1 after withdraw 500:", user1)

}
