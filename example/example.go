package example

import "log"

type MyLogger struct{}

func (m *MyLogger) Log(format string, v ...interface{}) {
	log.Printf(format, v...)
}

type Master struct {
	Logger    `inject:"logger"`
	Food      `inject:"example.Master.Food"`
	Transport `inject:"example.Master.Transport"`
}

type Farmer struct {
	Logger  `inject:"logger"`
	Machine `inject:"example.TillageMachine.Machine"`
}

func (f *Farmer) GetRice() {
	err := f.Machine.Run(3)
	if err != nil {
		f.Log("Machine breaks, no rice")
	}
	f.Log("Got rice")
}

type TillageMachine struct {
	Logger `inject:"logger"`
}

func (tm *TillageMachine) Run(n int) error {
	tm.Log("Tillage %d hours", n)
	return nil
}

type Driver struct {
	Logger `inject:"logger"`
	plane  string
}

func (d *Driver) Setup() error {
	d.plane = "Boeing787"
	return nil
}

func (d *Driver) Fly(src, dst string) {
	d.Log("%s Fly from %s to %s", d.plane, src, dst)
}
