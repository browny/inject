package inject

import (
	"reflect"
	"testing"

	"github.com/browny/inject/example"
	"github.com/stretchr/testify/suite"
)

func TestInjectTestSuite(t *testing.T) {
	suite.Run(t, new(InjectTestSuite))
}

type InjectTestSuite struct {
	suite.Suite
}

func (s *InjectTestSuite) TestWeave() {
	driver := example.Driver{}
	farmer := example.Farmer{}
	master := example.Master{}
	myLogger := example.MyLogger{}
	tillageMachine := example.TillageMachine{}

	depMap := map[interface{}][]string{
		&myLogger: []string{
			"logger",
		},
		&driver: []string{
			"example.Master.Transport",
		},
		&farmer: []string{
			"example.Master.Food",
		},
		&tillageMachine: []string{
			"example.TillageMachine.Machine",
		},
		&master: []string{},
	}

	graph, err := Weave(depMap)
	s.NoError(err)

	master.Food.GetRice()
	master.Transport.Fly("C++", "Go")

	f := graph[reflect.TypeOf(&example.Farmer{})].(*example.Farmer)
	f.Machine.Run(5)
}
