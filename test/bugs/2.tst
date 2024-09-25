# Bug 2: error checking task args for mixed lists.
# Type is checked on the first element of the list only.
# Fixed on 2018-03-29.
# This build should fail with message: "field '$' must be of type '[]string'"

default: bug-2

targets:

  bug-2:
    steps:
    - $: ['ls', 1]
