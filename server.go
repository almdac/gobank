package main

import (
	"net/http"

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
	}

	Withdraw struct {
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

	return c.JSON(http.StatusOK, acc)
}

func (b *Bank) Withdraw(c echo.Context) error {
	withdrawn := new(Withdraw)

	if err := c.Bind(withdrawn); err != nil {
		return err
	}
	if err := b.Authenticate(withdrawn.Email, withdrawn.Pass); err != nil {
		return err
	}

	acc := b.Accounts[withdrawn.Email]
	acc.Balance -= withdrawn.Value

	return c.JSON(http.StatusOK, acc)
}

func main() {
	e := echo.New()

	bank := &Bank{
		Accounts: make(map[string]*Account),
	}

	e.POST("/acc/create", bank.CreateAccount)
	e.PUT("/acc/withdraw", bank.Withdraw)
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
