package errs_test

import (
	"errors"
	"fmt"

	"github.com/fengjx/go-halo/errs"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// github.com/fengjx/go-halo/errs_test.ExampleNew_printf
	//         /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:17
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/fengjx/go-halo/errs/_test/_testmain.go:106
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func ExampleWithMessage() {
	cause := errors.New("whoops")
	err := errs.WithMessage(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func ExampleWithStack() {
	cause := errors.New("whoops")
	err := errs.WithStack(cause)
	fmt.Println(err)

	// Output: whoops
}

func ExampleWithStack_printf() {
	cause := errors.New("whoops")
	err := errs.WithStack(cause)
	fmt.Printf("%+v", err)

	// Example Output:
	// whoops
	// github.com/fengjx/go-halo/errs_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/fengjx/go-halo/errs/example_test.go:55
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/fengjx/go-halo/errs/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
	// github.com/fengjx/go-halo/errs_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/fengjx/go-halo/errs/example_test.go:56
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/fengjx/go-halo/errs/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errs.Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func fn() error {
	e1 := errors.New("error")
	e2 := errs.Wrap(e1, "inner")
	e3 := errs.Wrap(e2, "middle")
	return errs.Wrap(e3, "outer")
}

func ExampleCause() {
	err := fn()
	fmt.Println(err)
	fmt.Println(errs.Cause(err))

	// Output: outer: middle: inner: error
	// error
}

func ExampleWrap_extended() {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Example output:
	// error
	// github.com/fengjx/go-halo/errs_test.fn
	//         /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:47
	// github.com/fengjx/go-halo/errs_test.ExampleCause_printf
	//         /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:63
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/fengjx/go-halo/errs/_test/_testmain.go:104
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
	// github.com/fengjx/go-halo/errs_test.fn
	// 	  /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:48: inner
	// github.com/fengjx/go-halo/errs_test.fn
	//        /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:49: middle
	// github.com/fengjx/go-halo/errs_test.fn
	//      /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:50: outer
}

func ExampleWrapf() {
	cause := errors.New("whoops")
	err := errs.Wrapf(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Output: oh noes #2: whoops
}

func ExampleErrorf_extended() {
	err := fmt.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops: foo
	// github.com/fengjx/go-halo/errs_test.ExampleErrorf
	//         /home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:101
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/fengjx/go-halo/errs/_test/_testmain.go:102
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func Example_stackTrace() {
	type stackTracer interface {
		StackTrace() errs.StackTrace
	}

	err, ok := errs.Cause(fn()).(stackTracer)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	st := err.StackTrace()
	fmt.Printf("%+v", st[0:2]) // top two frames

	// Example output:
	// github.com/fengjx/go-halo/errs_test.fn
	//	/home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:47
	// github.com/fengjx/go-halo/errs_test.Example_stackTrace
	//	/home/dfc/src/github.com/fengjx/go-halo/errs/example_test.go:127
}

func ExampleCause_printf() {
	err := errs.Wrap(func() error {
		return func() error {
			return errors.New("hello world")
		}()
	}(), "failed")

	fmt.Printf("%v", err)

	// Output: failed: hello world
}