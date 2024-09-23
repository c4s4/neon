# This build file should fail: test2 -> test1

default: test1

targets:

  test1:
    depends: [test2]
    steps:
    - print: "test1"

  test2:
    steps:
    - print: "test2"
    - call: test1
