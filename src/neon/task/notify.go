package task

import (
	"github.com/gen2brain/beeep"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "notify",
		Func: notify,
		Args: reflect.TypeOf(notifyArgs{}),
		Help: `Desktop notification.

Arguments:

- notify: the title of the notification
- text: the notification text (optional)
- image: path to the notification image (optional)

Examples:

    # print a warning
    - notify: Warning
      text: This is a warning!`,
	})
}

type notifyArgs struct {
	Notify string
	Text   string `neon:"optional"`
	Image  string `neon:"file,optional"`
}

func notify(context *build.Context, args interface{}) error {
	params := args.(notifyArgs)
	return beeep.Notify(params.Notify, params.Text, params.Image)
}
