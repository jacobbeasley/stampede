package models

import (
	"log"
	"os"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/suite/v4"
)

func init() {
	os.Setenv("GO_ENV", "test")
	envy.Set("GO_ENV", "test")
	var err error
	DB, err = pop.Connect("test")
	if err != nil {
		log.Fatalf("Error re-initializing DB for models tests: %v", err)
	}
}

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(os.DirFS("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}
