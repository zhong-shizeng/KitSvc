package main

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/client"
	"github.com/TeaMeow/KitSvc/model"
	"github.com/TeaMeow/KitSvc/shared/token"
	"github.com/TeaMeow/KitSvc/version"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

//
var started = make(chan bool)
var c = client.NewClient("http://127.0.0.1:8080")
var ct client.Client

func TestMain(t *testing.T) {
	app := cli.NewApp()
	app.Name = "service"
	app.Version = version.Version
	app.Usage = "starts the service daemon."
	app.Action = func(c *cli.Context) {
		server(c, started)
	}
	app.Flags = serverFlags

	go app.Run(os.Args)

	<-started
}

func printErrors(e []error) {
	if len(e) != 0 {
		for _, v := range e {
			logrus.Error(v.Error())
		}
	}
}

func TestPostUser(t *testing.T) {
	assert := assert.New(t)

	u, errs := c.PostUser(&model.User{
		Username: "admin",
		Password: "testtest",
	})
	printErrors(errs)

	err := u.Compare("testtest")
	assert.True(err == nil)
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	u, errs := c.GetUser("admin")
	printErrors(errs)

	err := u.Compare("testtest")
	assert.True(err == nil)
}

func TestPostAuth(t *testing.T) {
	assert := assert.New(t)

	tkn, errs := c.PostAuth(&model.User{
		Username: "admin",
		Password: "testtest",
	})
	printErrors(errs)

	ctx, _ := token.Parse(tkn.Token, "4Rtg8BPKwixXy2ktDPxoMMAhRzmo9mmuZjvKONGPZZQSaJWNLijxR42qRgq0iBb5")
	assert.Equal(&token.Content{
		ID:       1,
		Username: "admin",
	}, ctx, "They should be equal.")

	ct = client.NewClientToken("http://127.0.0.1:8080", tkn.Token)
}

func TestPutUser(t *testing.T) {
	assert := assert.New(t)

	u, errs := ct.PutUser(1, &model.User{
		Username: "admin",
		Password: "newpassword",
	})
	printErrors(errs)

	err := u.Compare("newpassword")
	assert.True(err == nil, "They should be match.")
}

func TestDeleteUser(t *testing.T) {
	//assert := assert.New(t)
}