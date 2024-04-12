package test

import (
	"banners/app"
	"os"
	"testing"
)

var a app.App

func TestMain(m *testing.M) {

	a.Initialize()
	code := m.Run()
	os.Exit(code)
}
