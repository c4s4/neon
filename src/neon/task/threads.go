package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"sync"
	"strconv"
)

func init() {
	build.TaskMap["threads"] = build.TaskDescriptor{
		Constructor: Threads,
		Help: `Run steps in threads.

Arguments:

- threads: the number of threads to run. You can set it to _NCPU for the number
  of CPUs.
- data: a list filled with values to pass to threads in _data property.
- steps: the steps to run in threads.

Examples:

    # compute squares of 10 first integers in threads
    - threads: _NCPU
      data:    [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
      steps:
	  - 'square = _data * _data'
	  - print: '#{_data}^2 = #{square}'`,
	}
}

func Threads(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"threads", "data", "steps"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	threads, err := args.GetString("threads")
	if err != nil {
		threadInt, err := args.GetInteger("threads")
		if err != nil {
			return nil, fmt.Errorf("'threads' field must be an integer or an expression")
		}
		threads = strconv.Itoa(threadInt)
	}
	data, err := args.GetList("data")
	if err != nil {
		return nil, fmt.Errorf("'data' field of 'threads' must be a list")
	}
	steps, err := ParseSteps(target, args, "steps")
	if err != nil {
		return nil, err
	}
	return func(context *build.Context) error {
		_data := make(chan interface{}, len(data))
		for _, _d := range data {
			_data <- _d
		}
		_threads, _err := context.VM.EvaluateExpression(threads)
		if _err != nil {
			return fmt.Errorf("evaluating 'threads' field: %v", _err)
		}
		_nbThreads, _ok := _threads.(int)
		if !_ok {
			return fmt.Errorf("'threads' field must result as an integer")
		}
		_error := make(chan error, _nbThreads)
		var _waitGroup sync.WaitGroup
		_waitGroup.Add(_nbThreads)
		for _i := 0; _i < _nbThreads; _i++ {
			go RunThread(steps, context, _i, _data, &_waitGroup, _error)
		}
		_waitGroup.Wait()
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

func RunThread(steps []build.Step, ctx *build.Context, index int, data chan interface{}, wg *sync.WaitGroup, errors chan error) {
	ctx.Message("Thread %d started", index)
	defer ctx.Message("Thread %d done", index)
	defer wg.Done()
	stop := false
	for stop == false {
		select {
		case datum, ok := <-data:
			if ok {
				ctx.Message("Thread %d run with param: %v", index, datum)
				threadContext := ctx.Copy(index, datum)
				err := RunSteps(steps, threadContext)
				if err != nil {
					errors <- err
					stop = true
				}
			} else {
				stop = true
			}
		default:
			stop = true
		}
	}
}
