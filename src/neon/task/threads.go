package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"sync"
)

func init() {
	build.TaskMap["threads"] = build.TaskDescriptor{
		Constructor: Threads,
		Help: `Run steps in threads.

Arguments:

- threads: the number of threads to run. You can set it to _NCPU for the number
  of CPUs.
- input: a list filled with values to pass to threads in _data property.
- steps: the steps to run in threads.

Note:

This task sets two properties :
- _thread with the thread number (starting with 0)
- _input with the input for each thread.

If threads must output something, it should write it in _output. After threads
are done, _output will contain a list of all the outputs of threads.

Examples:

    # compute squares of 10 first integers in threads and put them in _output
    - threads: _NCPU
      input:   [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
      steps:
	  - 'square = _data * _data'
	  - '_output <- square'
	  - print: '#{_data}^2 = #{square}'
	# print squares on the console
	- print: '#{_output}'`,
	}
}

func Threads(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"threads", "input", "steps"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	var threads int
	var threadExpression string
	threads, err := args.GetInteger("threads")
	if err != nil {
		threadExpression, err = args.GetString("threads")
		if err != nil {
			return nil, fmt.Errorf("'threads' field must be an integer or an expression")
		}
	}
	var input []interface{}
	var inputExpression string
	input, err = args.GetList("input")
	if err != nil {
		inputExpression, err = args.GetString("input")
		if err != nil {
			return nil, fmt.Errorf("'input' field of 'threads' must be a list or an expression returning a list")
		}
	}
	steps, err := ParseSteps(target, args, "steps")
	if err != nil {
		return nil, err
	}
	return func(context *build.Context) error {
		if input == nil {
			_result, _err := context.EvaluateExpression(inputExpression)
			if err != nil {
				return fmt.Errorf("evaluating 'input' field: %v", _err)
			}
			var _ok bool
			input, _ok = _result.([]interface{})
			if !_ok {
				return fmt.Errorf("expression in 'input' field must return a list")
			}
		}
		_input := make(chan interface{}, len(input))
		for _, _d := range input {
			_input <- _d
		}
		if threadExpression != "" {
			_threads, _err := context.EvaluateExpression(threadExpression)
			if _err != nil {
				return fmt.Errorf("evaluating 'threads' field: %v", _err)
			}
			switch _t := _threads.(type) {
			case int:
				threads = _t
			case int64:
				threads = int(_t)
			default:
				return fmt.Errorf("'threads' field must result as an integer")
			}
		}
		_error := make(chan error, threads)
		var _waitGroup sync.WaitGroup
		_waitGroup.Add(threads)
		context.Message("Starting %d threads", threads)
		_output := make(chan interface{}, len(input))
		for _i := 0; _i < threads; _i++ {
			go RunThread(steps, context, _i, _input, _output, &_waitGroup, _error)
		}
		_waitGroup.Wait()
		var _out []interface{}
		stop := false
		for !stop {
			select {
			case o, ok := <-_output:
				if ok {
					_out = append(_out, o)
				} else {
					stop = true
				}
			default:
				stop = true
			}
		}
		context.SetProperty("_output", _out)
		select {
		case e, ok := <-_error:
			if ok {
				return e
			} else {
				return nil
			}
		default:
			return nil
		}
	}, nil
}

func RunThread(steps []build.Step, ctx *build.Context, index int, input chan interface{}, output chan interface{},
	wg *sync.WaitGroup, errors chan error) {
	ctx.Message("Thread %d started", index)
	defer ctx.Message("Thread %d done", index)
	defer wg.Done()
	for {
		select {
		case arg, ok := <-input:
			if ok {
				threadContext := ctx.NewThreadContext(index, arg, output)
				threadContext.Message("Thread %d iteration with data '%v'", index, arg)
				err := threadContext.Run(steps)
				out, _ := threadContext.GetProperty("_output")
				if err != nil {
					errors <- err
					return
				}
				if out != nil {
					output <- out
				}
			} else {
				return
			}
		default:
			return
		}
	}
}
