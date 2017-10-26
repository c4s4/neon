package build

import "strconv"

// ThreadContext is a context for a thread
type ThreadContext struct {
	Context
	Thread int
}

// NewThreadContext builds a context in a thread
func NewThreadContext(context *Context, thread int, data interface{}) *ThreadContext {
	threadContext := ThreadContext{
		Context: *context.Copy(),
		Thread:  thread,
	}
	threadContext.SetProperty("_thread", thread)
	threadContext.SetProperty("_data", data)
	return &threadContext
}

// Message print a message on the console
func (context *ThreadContext) Message(text string, args ...interface{}) {
	text = strconv.Itoa(context.Thread) + "| " + text
	Message(text, args...)
}
