package main

type Target struct {
	Depends []string
	Steps   []Step
	Name    string
}

func (t Target) Run() {
	for _, d := range t.Depends {
		target := build.Target(d)
		target.Run()
	}
	for _, s := range t.Steps {
		s.Run()
	}
}
