package task

import (
	"neon/build"
	"reflect"
	"sync"
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

Don't change current directory in threads as it would affect other threads as
well.`,
	})
}

type threadsArgs struct {
	Threads int
	Input   []interface{} `optional`
	Steps   build.Steps   `steps`
	Verbose bool          `optional`
}

func threads(context *build.Context, args interface{}) error {
	params := args.(threadsArgs)
	input := make(chan interface{}, len(params.Input))
	for _, d := range params.Input {
		input <- d
	}
	error := make(chan error, params.Threads)
	var wg sync.WaitGroup
	wg.Add(params.Threads)
	if params.Verbose {
		context.Message("Starting %d threads", params.Threads)
	}
	output := make(chan interface{}, len(input))
	for i := 0; i < params.Threads; i++ {
		go runThread(params.Steps, context, i, input, output, &wg, error, params.Verbose)
	}
	wg.Wait()
	var out []interface{}
	stop := false
	for !stop {
		select {
		case o, ok := <-output:
			if ok {
				out = append(out, o)
			} else {
				stop = true
			}
		default:
			stop = true
		}
	}
	context.SetProperty("_output", out)
	select {
	case e, ok := <-error:
		if ok {
			return e
		}
		return nil
	default:
		return nil
	}
}

func runThread(steps build.Steps, ctx *build.Context, index int, input chan interface{}, output chan interface{},
	wg *sync.WaitGroup, errors chan error, verbose bool) {
	if verbose {
		ctx.Message("Thread %d started", index)
		defer ctx.Message("Thread %d done", index)
	}
	defer wg.Done()
	for {
		select {
		case arg, ok := <-input:
			if ok {
				threadContext := ctx.NewThreadContext(index, arg, output)
				if verbose {
					threadContext.Message("Thread %d iteration with input '%v'", index, arg)
				}
				err := steps.Run(threadContext)
				out, _ := threadContext.GetProperty("_output")
				if err != nil {
					errors <- err
					return
				}
				if out != nil {
					output <- out
					if verbose {
						threadContext.Message("Thread %d output '%v'", index, out)
					}
				}
			} else {
				return
			}
		default:
			return
		}
	}
}
