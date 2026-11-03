package task

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/c4s4/neon/neon/build"
)

const (
	propertyThread = "_thread"
	propertyInput  = "_input"
	propertyOutput = "_output"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "threads",
		Func: threads,
		Args: reflect.TypeOf(threadsArgs{}),
		Help: `Run steps in threads.

Arguments:

- threads: number of threads to run (integer).
- input: values to pass to threads in _input property (list, optional).
- steps: steps to run in threads (steps).
- verbose: if you want thread information on console, defaults to false
  (boolean, optional).

Examples:

    # compute squares of 10 first integers in threads and put them in _output
    - threads: =_NCPU
      input:   =range(10)
      steps:
      - '_output = _input * _input'
      - print: '#{_input}^2 = #{_output}'
    # print squares on the console
    - print: '#{_output}'

Notes:

- You might set number of threads to '_NCPU' which is the number of cores in
  the CPU of the machine.
- Property _thread is set with the thread number (starting with 0)
- Property _input is set with the input for each thread.
- Property _output is set with the output of the threads.
- Each thread should write its output in property _output.

Context of the build is cloned in each thread so that you can read and write
properties, they won't affect other threads. But all properties will be lost
when thread is done, except for _output that will be appended to other in
resulting _output property.

Current directory is shared by all threads: changing to another directory
while threads are running (with chdir or neon tasks, or calling a target of a
build in another directory) raises an error. Number of threads must be at
least 1. When a thread fails, other threads stop and the first error is
returned.`,
	})
}

type threadsArgs struct {
	Threads int
	Input   []interface{} `neon:"optional"`
	Steps   build.Steps   `neon:"steps"`
	Verbose bool          `neon:"optional"`
}

func threads(context *build.Context, args interface{}) error {
	params := args.(threadsArgs)
	if params.Threads < 1 {
		return fmt.Errorf("number of threads must be at least 1 (%d provided)", params.Threads)
	}
	if params.Input == nil {
		params.Input = make([]interface{}, params.Threads)
	}
	input := make(chan interface{}, len(params.Input))
	for _, d := range params.Input {
		input <- d
	}
	close(input)
	build.EnterParallel()
	defer build.LeaveParallel()
	errs := make(chan error, params.Threads)
	// stop is closed on first error to stop other threads
	stop := make(chan struct{})
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(params.Threads)
	if params.Verbose {
		context.MessageArgs("Starting %d threads", params.Threads)
	}
	output := make(chan interface{}, len(params.Input))
	for i := 0; i < params.Threads; i++ {
		go func(index int) {
			defer wg.Done()
			if err := runThread(params.Steps, context, index, input, output, stop, params.Verbose); err != nil {
				errs <- err
				once.Do(func() { close(stop) })
			}
		}(i)
	}
	wg.Wait()
	close(output)
	close(errs)
	var out []interface{}
	for o := range output {
		out = append(out, o)
	}
	context.SetProperty(propertyOutput, out)
	// return first error, if any
	return <-errs
}

func runThread(steps build.Steps, ctx *build.Context, index int, input <-chan interface{}, output chan<- interface{},
	stop <-chan struct{}, verbose bool) error {
	if verbose {
		ctx.MessageArgs("Thread %d started", index)
		defer ctx.MessageArgs("Thread %d done", index)
	}
	for arg := range input {
		select {
		case <-stop:
			return nil
		default:
		}
		threadContext := ctx.Copy()
		threadContext.SetProperty(propertyThread, index)
		threadContext.SetProperty(propertyInput, arg)
		threadContext.DelProperty(propertyOutput)
		if verbose {
			threadContext.MessageArgs("Thread %d iteration with input '%v'", index, arg)
		}
		if err := steps.Run(threadContext); err != nil {
			return err
		}
		out, _ := threadContext.GetProperty(propertyOutput)
		if out != nil {
			output <- out
			if verbose {
				threadContext.MessageArgs("Thread %d output '%v'", index, out)
			}
		}
	}
	return nil
}
