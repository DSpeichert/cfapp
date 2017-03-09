package main

import (
	"bytes"
	"encoding/base64"
	"math/rand"

	"encoding/json"
	"net/http"

	"fmt"

	"golang.org/x/crypto/scrypt"
)

// swagger:model
type Customer struct {
	ID           int           `gorm:"primary_key" json:"id" form:"id"`
	Name         string        `json:"name" form:"name"`
	Email        string        `json:"email" form:"email"`
	Password     string        `json:"password" form:"password"`
	Salt         string        `json:"salt" form:"salt"`
	Certificates []Certificate `json:"certificates" form:"certificates"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (c *Customer) makeSalt() {
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	c.Salt = string(b)
}

func (c *Customer) HashPassword() {
	c.makeSalt()
	dk, _ := scrypt.Key([]byte(c.Password), []byte(c.Salt), 16384, 8, 1, 32)
	c.Password = base64.StdEncoding.EncodeToString(dk)
}

// swagger:model
type Certificate struct {
	ID          int    `gorm:"primary_key" json:"id" form:"id"`
	CustomerID  int    `json:"customer_id" form:"customer_id"`
	Active      bool   `json:"active" form:"active"`
	Certificate string `gorm:"type:text" json:"certificate" form:"certificate"`
	Key         string `gorm:"type:text" json:"key" form:"key"`
}

func (c *Certificate) PostUpdateHook() {
	if notify_url != "" {
		n := CertificateNotification{CertificateID: c.ID, Active: c.Active}
		j, _ := json.MarshalIndent(n, "", `    `)
		_, err := http.Post(notify_url, "application/json", bytes.NewBuffer(j))
		if err != nil {
			fmt.Errorf("%s\n", err.Error())
		}
	}
}

type QueryResponse struct {
	Total  int         `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Items  interface{} `json:"items"`
}

type CertificateNotification struct {
	CertificateID int  `json:"id"`
	Active        bool `json:"active"`
}
