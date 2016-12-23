package main

type Target struct {
	Name    string
	Doc     string
	Build   *Build
	Depends []string
	Steps   []string
}

func (target *Target) Init(build *Build, name string) {
	target.Build = build
	target.Name = name
}

func (target *Target) Run() {
	for _, depend := range target.Depends {
		dependency := target.Build.Target(depend)
		dependency.Run()
	}
	PrintTarget("Running target " + target.Name)
	for _, step := range target.Steps {
		target.Build.Context.Execute(step)
	}
}
