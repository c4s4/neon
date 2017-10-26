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
- data: a list filled with values to pass to threads in _data property.
- steps: the steps to run in threads.

Note:

This task sets two properties :
- _data with the data for each thread.
- _thread with the thread number (starting with 0)

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
	var threads int
	var threadExpression string
	threads, err := args.GetInteger("threads")
	if err != nil {
		threadExpression, err = args.GetString("threads")
		if err != nil {
			return nil, fmt.Errorf("'threads' field must be an integer or an expression")
		}
	}
	var data []interface{}
	var dataExpression string
	data, err = args.GetList("data")
	if err != nil {
		dataExpression, err = args.GetString("data")
		if err != nil {
			return nil, fmt.Errorf("'data' field of 'threads' must be a list or an expression returning a list")
		}
	}
	steps, err := ParseSteps(target, args, "steps")
	if err != nil {
		return nil, err
	}
	return func(context *build.Context) error {
		if data == nil {
			_result, _err := context.EvaluateExpression(dataExpression)
			if err != nil {
				return fmt.Errorf("evaluating 'data' field: %v", _err)
			}
			var _ok bool
			data, _ok = _result.([]interface{})
			if !_ok {
				return fmt.Errorf("expression in 'data' field must return a list")
			}
		}
		_data := make(chan interface{}, len(data))
		for _, _d := range data {
			_data <- _d
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
		for _i := 0; _i < threads; _i++ {
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
	for {
		select {
		case datum, ok := <-data:
			if ok {
				ctx.Message("Thread %d run with param: %v", index, datum)
				threadContext := build.NewThreadContext(ctx, index, datum)
				err := threadContext.Context.Run(steps)
				if err != nil {
					errors <- err
					return
				}
			} else {
				return
			}
		default:
			return
		}
	}
}
