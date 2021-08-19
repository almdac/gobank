package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type (
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	Account struct {
		User    User    `json:"user"`
		Pass    string  `json:"pass"`
		Balance float64 `json:"balance"`
	}

	Bank struct {
		Accounts map[string]*Account
		Mutexes  map[string]*sync.Mutex
	}

	Value struct {
		Email string  `json:"email"`
		Pass  string  `json:"pass"`
		Value float64 `json:"value"`
	}
)

func (b *Bank) CreateAccount(c echo.Context) error {
	acc := new(Account)

	if err := c.Bind(acc); err != nil {
		return err
	}

	b.Accounts[acc.User.Email] = acc
	b.Mutexes[acc.User.Email] = new(sync.Mutex)

	return c.JSON(http.StatusOK, acc)
}

func (b *Bank) Withdraw(c echo.Context) error {
	withdrawn := new(Value)

	if err := c.Bind(withdrawn); err != nil {
		return err
	}
	if err := b.Authenticate(withdrawn.Email, withdrawn.Pass); err != nil {
		return err
	}

	acc := b.Accounts[withdrawn.Email]
	mutex := b.Mutexes[withdrawn.Email]

	mutex.Lock()
	acc.Balance -= withdrawn.Value
	defer mutex.Unlock()

	return c.JSON(http.StatusOK, acc)
}

func (b *Bank) Deposit(c echo.Context) error {
	deposit := new(Value)

	if err := c.Bind(deposit); err != nil {
		return err
	}
	if err := b.Authenticate(deposit.Email, deposit.Pass); err != nil {
		return err
	}

	acc := b.Accounts[deposit.Email]
	mutex := b.Mutexes[deposit.Email]

	mutex.Lock()
	acc.Balance += deposit.Value
	defer mutex.Unlock()

	return c.JSON(http.StatusOK, acc)
}

func main() {
	e := echo.New()

	bank := &Bank{
		Accounts: make(map[string]*Account),
		Mutexes:  make(map[string]*sync.Mutex),
	}

	e.POST("/acc/create", bank.CreateAccount)
	e.PUT("/acc/withdraw", bank.Withdraw)
	e.PUT("/acc/deposit", bank.Deposit)
	e.Logger.Fatal(e.Start(":3000"))
}

func (b *Bank) Authenticate(email, pass string) error {
	acc, ok := b.Accounts[email]

	if !ok {
		return echo.ErrNotFound
	}

	if acc.Pass != pass {
		return echo.ErrForbidden
	}

	return nil
}
