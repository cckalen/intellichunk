# goapp-template
A template repository for a skeleton go app with command line handling.

It has a `cmd` package with a `root` (mostly boilerplate) that names your application and defines options for all
commands (for example logging level).
There are two commands `hello` which is an umbrella for subcommands. There is `hello_world` which is a subcommand of `hello`.
Subcommands can be nested to any depth. Files are typically named after the path from the root.

The `internal` folder is special to go as the content is considered private to the module/application. Logic you expect to be used by other modules are considered the API of the module, and they should be in top level folders within the module.

The `internal/check` is an example of a simple enumeration check for valid loglevel names.
The `internal/logging` configures logging based on a loglevel - it is called from the `root` command during initialization.
The `internal/example` contains a `greet` file with various `Great` functions called from `hello_world` command.
The `example_test` is a test package that tests the (module internal) API/contract of the example logic. Always put tests in a separate package unless testing cannot be done without access to the implementatation of what is being tested. In that case only put those special test inside the package being tested. (Testlogic will still be dropped from reqular builds).
The `example_test/greet_test` tests the various `Greet` functions
The `example_test/testutils_sample_test` shows some simple examples of how to use the `testutils` module.

== Create your application from the template

To make this into your working go application:
* Initialize a github repo with the repo of this template, or copy everything over manually.
* Change the folder name "goapp-template" to the name of your application/package.
* Also change all import references containing "goapp-template" to the same application name.
* There are TODO marks where editing needs to be done.
    * For example changing "myapp" and "Myapp" in strings/text to the wanted name of your application.
* Add suitable commands for you app by changing the "hello" example in `cmd`, or createing new commands and deleting the hello
* The example package contain the functionality used by the example hello command line. There are also tests that show bare bones use of testutils package.
* The `Test_examples_of_testutils` shows basic usage of the testutils
* The example package can naturally be deleted from your app.

== Running the app
Run with `go run .` to run the `main` (which is boilerplate and you want to keep in your app). Try running with these additional arguments:
* hello
* hello world
* hello world blah
* hello world blah blah

And combine with the flags `-w` (wonderful) and `-u` (upper case).

For example:
```
go run . hello world albert -w
```

== Running Tests
Tests are run with `go test`, the argument `./...` runs all tests anywhere in the file structure, the `-v` outputs the result of each test, and `-count=1` forces go to build first and not run tests from cached earlier build (which is irritating when you changed a source file and did not do `go run` before running tests, and therefore does not test what you just changed).

```
go test ./... -count=1 -v
?       github.com/wyrth-io/goapp-template      [no test files]
?       github.com/wyrth-io/goapp-template/cmd  [no test files]
?       github.com/wyrth-io/goapp-template/internal/check       [no test files]
=== RUN   Test_Greet_returns_hello_without_name
--- PASS: Test_Greet_returns_hello_without_name (0.00s)
=== RUN   Test_Greet_returns_hello_with_name_when_given
--- PASS: Test_Greet_returns_hello_with_name_when_given (0.00s)
=== RUN   Test_GreetUpper_returns_HELLO_without_name
--- PASS: Test_GreetUpper_returns_HELLO_without_name (0.00s)
=== RUN   Test_GreetUpper_returns_HELLO_with_name_upcased_when_given
--- PASS: Test_GreetUpper_returns_HELLO_with_name_upcased_when_given (0.00s)
=== RUN   Test_GreetWonderful_returns_hello_world_without_name
--- PASS: Test_GreetWonderful_returns_hello_world_without_name (0.00s)
=== RUN   Test_GreetWonderful_returns_hello_with_name_when_given
--- PASS: Test_GreetWonderful_returns_hello_with_name_when_given (0.00s)
=== RUN   Test_GreetWonderfulUpper_returns_HELLO_WORLD_without_name
--- PASS: Test_GreetWonderfulUpper_returns_HELLO_WORLD_without_name (0.00s)
=== RUN   Test_GreetWorldUpper_returns_hello_with_name_when_given
--- PASS: Test_GreetWorldUpper_returns_hello_with_name_when_given (0.00s)
=== RUN   Test_examples_of_testutils
--- PASS: Test_examples_of_testutils (0.00s)
PASS
ok      github.com/wyrth-io/goapp-template/internal/example     0.333s
?       github.com/wyrth-io/goapp-template/internal/logging     [no test files]
```

### Code Quality
In order for code to be merged it must pass all tests and all configured linting. The linters check both hard and
stylistic issues. Any reported problem must be fixed before merging.

**Setup**</br>
* We build with latest go (currently 1.91) - if you do not have that installed go here  [https://go.dev.dl](https://go.dev.dl)
* We run lint with `golangci-lint` and you should install it locally so you can check your code. Follow instructions
  here [https://golangci-lint.run/usage/install/](https://golangci-lint.run/usage/install/)

**Continously, or at least before PR**</br>
* Run `go test ./... -count=1 -v` and make sure all tests are green / ok.
* Run `golangci-lint run` which runs linting on the entire project, make sure there are no warning or errors

## Github Actions / workflows
This template contains github actions workflow to run all tests and configured linters executed by `golangci-lint`. The workflow runs on push and pull requests to the origin repo.
