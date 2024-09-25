# Bug 4: Bad error message when parent target not found calling super
# Fixed on 2018-05-24
# Must fail with message: "target 'bug-4' not found in parent build files"

default: bug-4

targets:

  bug-4:
    steps:
    - super:
