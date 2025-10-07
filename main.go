package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
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
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += sum
}

func (u *User) Withdraw(sum float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Balance < sum {
		return errors.New("insufficient funds")
	}
	u.Balance -= sum
	return nil
}

func (ps *PaymentSystem) AddUser(user *User) {
	ps.Users[user.ID] = user
}

func (ps *PaymentSystem) AddTransaction(transaction Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, transaction)
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

func (ps *PaymentSystem) Worker(ch chan Transaction, wg *sync.WaitGroup) error {
	for t := range ch {
		if err := ps.ProcessingTransactions(t); err != nil {
			return err
		}
	}
	defer wg.Done()
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

	paymentSystem.AddTransaction(Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 200})
	paymentSystem.AddTransaction(Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 50})
	paymentSystem.AddTransaction(Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 100})
	paymentSystem.AddTransaction(Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 140})

	ch := make(chan Transaction, len(paymentSystem.TransactionQueue))
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go paymentSystem.Worker(ch, &wg)
	}

	for _, t := range paymentSystem.TransactionQueue {
		ch <- t
	}

	close(ch)
	wg.Wait()

	fmt.Println("\nTotal:")
	fmt.Println("user1:", user1)
	fmt.Println("user2:", user2)
}
