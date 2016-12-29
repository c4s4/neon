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

func (target *Target) Run() error {
	for _, depend := range target.Depends {
		dependency, err := target.Build.Target(depend)
		if err != nil {
			return err
		}
		dependency.Run()
	}
	PrintTarget("Running target " + target.Name)
	for _, step := range target.Steps {
		err := target.Build.Context.Execute(step)
		if err != nil {
			return err
		}
	}
	return nil
}
