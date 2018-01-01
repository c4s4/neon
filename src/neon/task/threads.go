package task

import (
	"neon/build"
	"sync"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "threads",
		Func: Threads,
		Args: reflect.TypeOf(ThreadsArgs{}),
		Help: `Run steps in threads.

Arguments:

- threads: the number of threads to run. You can set it to _NCPU for the number
  of CPUs.
- input: a list filled with values to pass to threads in _input property.
- steps: the steps to run in threads.
- verbose: tells if threads information should be printed on console (optional,
  boolean that defaults to false).

Note:

This task sets two properties :
- _thread with the thread number (starting with 0)
- _input with the input for each thread.

Context of the build is cloned in each thread so that you can read and write
properties, they won't affect other threads. But all properties will be lost
when thread is done.

If threads must output something, they must write it in _output property.
After threads are done, _output will contain a list of all the outputs of
threads.

Don't change current directory in threads as it would affect other threads as
well.

Examples:

    # compute squares of 10 first integers in threads and put them in _output
    - threads: _NCPU
      input:   [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
      steps:
      - '_output = _input * _input'
      - print: '#{_input}^2 = #{_output}'
    # print squares on the console
    - print: '#{_output}'`,
	})
}

type ThreadsArgs struct {
	Threads int
	Input   []interface{} `optional`
	Steps   []build.Step  `steps`
	Verbose bool          `optional`
}

func Threads(context *build.Context, args interface{}) error {
	params := args.(ThreadsArgs)
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
		go RunThread(params.Steps, context, i, input, output, &wg, error, params.Verbose)
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
		} else {
			return nil
		}
	default:
		return nil
	}
}

func RunThread(steps []build.Step, ctx *build.Context, index int, input chan interface{}, output chan interface{},
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
				err := threadContext.Run(steps)
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
