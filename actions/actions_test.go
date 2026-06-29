package actions

import (
	"log"
	"os"
	"testing"

	"buffalo-app/models"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/suite/v4"
)

func init() {
	os.Setenv("GO_ENV", "test")
	envy.Set("GO_ENV", "test")
	var err error
	models.DB, err = pop.Connect("test")
	if err != nil {
		log.Fatalf("Error re-initializing models.DB for actions tests: %v", err)
	}
}

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), os.DirFS("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}
