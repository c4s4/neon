default: test1

targets:

  test1:
    depends: [test2]
    steps:
    - print: "test1"

  test2:
    depends: [test1]
    steps:
    - print: "test2"
