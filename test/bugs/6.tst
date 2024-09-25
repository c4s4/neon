# Bug 6: Doesn't load context script of parents build files
# This build must succeed and print "bug 6"

context:
- ./6.ank

targets:

  bug:
    steps:
    - print: "bug 6"
