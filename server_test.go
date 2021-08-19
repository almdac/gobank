package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	e := echo.New()
	bank := &Bank{
		Accounts: map[string]*Account{},
	}

	// Get body
	body := &Account{
		User: User{
			Name:  "Almir Menezes da Cunha Júnior",
			Email: "almirmdacunha@gmail.com",
		},
		Pass:    "ut9na5eb",
		Balance: 50,
	}
	bodyBytes, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(bodyBytes)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodPost, "/acc/create", bodyReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Create echo context
	c := e.NewContext(req, rec)

	// Validate
	if assert.NoError(t, bank.CreateAccount(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate if account was created
		createdAcc := bank.Accounts[body.User.Email]
		assert.Equal(t, *body, *createdAcc)
	}

}

func TestWithdraw(t *testing.T) {
	e := echo.New()
	bank := &Bank{
		Accounts: map[string]*Account{
			"almirmdacunha@gmail.com": {
				User: User{
					Name:  "Almir Menezes da Cunha Júnior",
					Email: "almirmdacunha@gmail.com",
				},
				Pass:    "ut9na5eb",
				Balance: 50,
			},
		},
	}

	// Get body
	body := &Withdraw{
		Email: "almirmdacunha@gmail.com",
		Pass:  "ut9na5eb",
		Value: 20,
	}
	bodyBytes, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(bodyBytes)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodPut, "/acc/withdraw", bodyReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Create echo context
	c := e.NewContext(req, rec)

	// Validate
	acc := bank.Accounts[body.Email]
	oldBalance := acc.Balance
	if assert.NoError(t, bank.Withdraw(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate if balance was correctly modified
		assert.Equal(t, oldBalance-body.Value, acc.Balance)
	}
}
