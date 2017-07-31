# inject [![GoDoc](https://godoc.org/github.com/browny/inject?status.svg)](http://godoc.org/github.com/browny/inject)
the improved version of https://github.com/facebookgo/inject

## Usage
See how it works in `inject_test.go` by running `go test`

```go
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
```
