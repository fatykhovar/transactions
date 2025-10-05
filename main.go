package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []*Transaction
}

func (u *User) String() string {
	return fmt.Sprintf("Name: %s, Amount: %.2f", u.Name, u.Balance)
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

func (ps *PaymentSystem) AddUser(user *User) {
	ps.Users[user.ID] = user
}

func (ps *PaymentSystem) AddTransaction(transaction *Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, transaction)
}

func (ps *PaymentSystem) ProcessingTransactions() error {
	if len(ps.TransactionQueue) == 0 {
		return errors.New("no transactions")
	}
	fromID, ok := ps.Users[ps.TransactionQueue[0].FromID]
	if !ok {
		return errors.New("user FromID not found")
	}
	toID, ok := ps.Users[ps.TransactionQueue[0].ToID]
	if !ok {
		return errors.New("user ToID not found")
	}

	if err := fromID.Withdraw(ps.TransactionQueue[0].Amount); err != nil {
		return fmt.Errorf("transaction %v: %w", ps.TransactionQueue[0], err)
	}
	fmt.Printf("After withdraw: %.2f, from: %v\n", ps.TransactionQueue[0].Amount, fromID)

	toID.Deposit(ps.TransactionQueue[0].Amount)
	fmt.Printf("After deposit: %.2f, to: %v\n", ps.TransactionQueue[0].Amount, toID)

	ps.TransactionQueue = ps.TransactionQueue[1:]
	return nil
}

func main() {
	user1 := &User{ID: uuid.NewString(), Name: "John", Balance: 1000}
	user2 := &User{ID: uuid.NewString(), Name: "Linda", Balance: 500}
	fmt.Println("new user1:", user1)
	fmt.Println("new user2:", user2)
	fmt.Println()

	paymentSystem := PaymentSystem{Users: make(map[string]*User), TransactionQueue: make([]*Transaction, 0)}

	paymentSystem.AddUser(user1)
	paymentSystem.AddUser(user2)

	transaction1 := &Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 200}
	transaction2 := &Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 50}

	paymentSystem.AddTransaction(transaction1)
	paymentSystem.AddTransaction(transaction2)

	for len(paymentSystem.TransactionQueue) > 0 {
		if err := paymentSystem.ProcessingTransactions(); err != nil {
			fmt.Println(err)
			break
		}
	}
	fmt.Println("\nTotal:")
	fmt.Println("user1:", user1)
	fmt.Println("user2:", user2)
}
