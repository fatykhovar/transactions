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
	TransactionQueue []Transaction
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

func (ps *PaymentSystem) AddTransaction(t Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, t)
}

func (ps *PaymentSystem) ProcessingTransactions(t Transaction) error {
	fromID, ok := ps.Users[t.FromID]
	if !ok {
		return errors.New("user FromID not found")
	}
	toID, ok := ps.Users[t.ToID]
	if !ok {
		return errors.New("user ToID not found")
	}

	if err := fromID.Withdraw(t.Amount); err != nil {
		return fmt.Errorf("transaction %v: %w", t, err)
	}
	fmt.Printf("After withdraw: %.2f, from: %v\n", t.Amount, fromID)

	toID.Deposit(t.Amount)
	fmt.Printf("After deposit: %.2f, to: %v\n", t.Amount, toID)

	return nil
}

func main() {
	user1 := &User{ID: uuid.NewString(), Name: "John", Balance: 1000}
	user2 := &User{ID: uuid.NewString(), Name: "Linda", Balance: 500}
	fmt.Println("new user1:", user1)
	fmt.Println("new user2:", user2)
	fmt.Println()

	paymentSystem := PaymentSystem{Users: make(map[string]*User), TransactionQueue: make([]Transaction, 0)}

	paymentSystem.AddUser(user1)
	paymentSystem.AddUser(user2)

	transaction1 := Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 200}
	transaction2 := Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 50}

	paymentSystem.AddTransaction(transaction1)
	paymentSystem.AddTransaction(transaction2)

	for _, t := range paymentSystem.TransactionQueue {
		if err := paymentSystem.ProcessingTransactions(t); err != nil {
			fmt.Println(err)
			break
		}
	}
	fmt.Println("\nTotal:")
	fmt.Println("user1:", user1)
	fmt.Println("user2:", user2)
}
