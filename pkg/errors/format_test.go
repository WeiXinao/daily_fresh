package errors

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error",
	}, {
		New("error"),
		"%+v",
		"error\n" +
			"imooc/mxshop/pkg/errors.TestFormatNew\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:26",
	}, {
		New("error"),
		"%q",
		`"error"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"error\n" +
			"imooc/mxshop/pkg/errors.TestFormatErrorf\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:56",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrap(New("error"), "error2"),
		"%s",
		"error2",
	}, {
		Wrap(New("error"), "error2"),
		"%v",
		"error2",
	}, {
		Wrap(New("error"), "error2"),
		"%+v",
		"error\n" +
			"imooc/mxshop/pkg/errors.TestFormatWrap\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:82",
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error",
	}, {
		Wrap(io.EOF, "error"),
		"%v",
		"error",
	}, {
		Wrap(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"error\n" +
			"imooc/mxshop/pkg/errors.TestFormatWrap\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:96",
	}, {
		Wrap(Wrap(io.EOF, "error1"), "error2"),
		"%+v",
		"EOF\n" +
			"error1\n" +
			"imooc/mxshop/pkg/errors.TestFormatWrap\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:103\n",
	}, {
		Wrap(New("error with space"), "context"),
		"%q",
		`"context"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%v",
		"error2",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"EOF\n" +
			"error2\n" +
			"imooc/mxshop/pkg/errors.TestFormatWrapf\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:134",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"imooc/mxshop/pkg/errors.TestFormatWrapf\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:149",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWithStack(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		WithStack(io.EOF),
		"%s",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%v",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%+v",
		[]string{"EOF",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:175"},
	}, {
		WithStack(New("error")),
		"%s",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%v",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%+v",
		[]string{"error",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:189",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:189"},
	}, {
		WithStack(WithStack(io.EOF)),
		"%+v",
		[]string{"EOF",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:197",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:197"},
	}, {
		WithStack(WithStack(Wrapf(io.EOF, "message"))),
		"%+v",
		[]string{"EOF",
			"message",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:205",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:205",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:205"},
	}, {
		WithStack(Errorf("error%d", 1)),
		"%+v",
		[]string{"error1",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:216",
			"imooc/mxshop/pkg/errors.TestFormatWithStack\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:216"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatWithMessage(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		WithMessage(New("error"), "error2"),
		"%s",
		[]string{"error2"},
	}, {
		WithMessage(New("error"), "error2"),
		"%v",
		[]string{"error2"},
	}, {
		WithMessage(New("error"), "error2"),
		"%+v",
		[]string{
			"error",
			"imooc/mxshop/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:244",
			"error2"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%s",
		[]string{"addition1"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%v",
		[]string{"addition1"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%+v",
		[]string{"EOF", "addition1"},
	}, {
		WithMessage(WithMessage(io.EOF, "addition1"), "addition2"),
		"%v",
		[]string{"addition2"},
	}, {
		WithMessage(WithMessage(io.EOF, "addition1"), "addition2"),
		"%+v",
		[]string{"EOF", "addition1", "addition2"},
	}, {
		Wrap(WithMessage(io.EOF, "error1"), "error2"),
		"%+v",
		[]string{"EOF", "error1", "error2",
			"imooc/mxshop/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:272"},
	}, {
		WithMessage(Errorf("error%d", 1), "error2"),
		"%+v",
		[]string{"error1",
			"imooc/mxshop/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:278",
			"error2"},
	}, {
		WithMessage(WithStack(io.EOF), "error"),
		"%+v",
		[]string{
			"EOF",
			"imooc/mxshop/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:285",
			"error"},
	}, {
		WithMessage(Wrap(WithStack(io.EOF), "inside-error"), "outside-error"),
		"%+v",
		[]string{
			"EOF",
			"imooc/mxshop/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:293",
			"inside-error",
			"imooc/mxshop/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:293",
			"outside-error"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatGeneric(t *testing.T) {
	starts := []struct {
		err  error
		want []string
	}{
		{New("new-error"), []string{
			"new-error",
			"imooc/mxshop/pkg/errors.TestFormatGeneric\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:315"},
		}, {Errorf("errorf-error"), []string{
			"errorf-error",
			"imooc/mxshop/pkg/errors.TestFormatGeneric\n" +
				"\t.+/imooc/mxshop/pkg/errors/format_test.go:319"},
		}, {errors.New("errors-new-error"), []string{
			"errors-new-error"},
		},
	}

	wrappers := []wrapper{
		{
			func(err error) error { return WithMessage(err, "with-message") },
			[]string{"with-message"},
		}, {
			func(err error) error { return WithStack(err) },
			[]string{
				"imooc/mxshop/pkg/errors.(func·002|TestFormatGeneric.func2)\n\t" +
					".+/imooc/mxshop/pkg/errors/format_test.go:333",
			},
		}, {
			func(err error) error { return Wrap(err, "wrap-error") },
			[]string{
				"wrap-error",
				"imooc/mxshop/pkg/errors.(func·003|TestFormatGeneric.func3)\n\t" +
					".+/imooc/mxshop/pkg/errors/format_test.go:339",
			},
		}, {
			func(err error) error { return Wrapf(err, "wrapf-error%d", 1) },
			[]string{
				"wrapf-error1",
				"imooc/mxshop/pkg/errors.(func·004|TestFormatGeneric.func4)\n\t" +
					".+/imooc/mxshop/pkg/errors/format_test.go:346",
			},
		},
	}

	for s := range starts {
		err := starts[s].err
		want := starts[s].want
		testFormatCompleteCompare(t, s, err, "%+v", want, false)
		testGenericRecursive(t, err, want, wrappers, 3)
	}
}

func wrappedNew(message string) error { // This function will be mid-stack inlined in go 1.12+
	return New(message)
}

func TestFormatWrappedNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+v",
		"error\n" +
			"imooc/mxshop/pkg/errors.wrappedNew\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:364\n" +
			"imooc/mxshop/pkg/errors.TestFormatWrappedNew\n" +
			"\t.+/imooc/mxshop/pkg/errors/format_test.go:373",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}

var stackLineR = regexp.MustCompile(`\.`)

// parseBlocks parses input into a slice, where:
//  - incase entry contains a newline, its a stacktrace
//  - incase entry contains no newline, its a solo line.
//
// Detecting stack boundaries only works incase the WithStack-calls are
// to be found on the same line, thats why it is optionally here.
//
// Example use:
//
// for _, e := range blocks {
//   if strings.ContainsAny(e, "\n") {
//     // Match as stack
//   } else {
//     // Match as line
//   }
// }
//
func parseBlocks(input string, detectStackboundaries bool) ([]string, error) {
	var blocks []string

	stack := ""
	wasStack := false
	lines := map[string]bool{} // already found lines

	for _, l := range strings.Split(input, "\n") {
		isStackLine := stackLineR.MatchString(l)

		switch {
		case !isStackLine && wasStack:
			blocks = append(blocks, stack, l)
			stack = ""
			lines = map[string]bool{}
		case isStackLine:
			if wasStack {
				// Detecting two stacks after another, possible cause lines match in
				// our tests due to WithStack(WithStack(io.EOF)) on same line.
				if detectStackboundaries {
					if lines[l] {
						if len(stack) == 0 {
							return nil, errors.New("len of block must not be zero here")
						}

						blocks = append(blocks, stack)
						stack = l
						lines = map[string]bool{l: true}
						continue
					}
				}

				stack = stack + "\n" + l
			} else {
				stack = l
			}
			lines[l] = true
		case !isStackLine && !wasStack:
			blocks = append(blocks, l)
		default:
			return nil, errors.New("must not happen")
		}

		wasStack = isStackLine
	}

	// Use up stack
	if stack != "" {
		blocks = append(blocks, stack)
	}
	return blocks, nil
}

func testFormatCompleteCompare(t *testing.T, n int, arg interface{}, format string, want []string, detectStackBoundaries bool) {
	gotStr := fmt.Sprintf(format, arg)

	got, err := parseBlocks(gotStr, detectStackBoundaries)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("test %d: fmt.Sprintf(%s, err) -> wrong number of blocks: got(%d) want(%d)\n got: %s\nwant: %s\ngotStr: %q",
			n+1, format, len(got), len(want), prettyBlocks(got), prettyBlocks(want), gotStr)
	}

	for i := range got {
		if strings.ContainsAny(want[i], "\n") {
			// Match as stack
			match, err := regexp.MatchString(want[i], got[i])
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("test %d: block %d: fmt.Sprintf(%q, err):\ngot:\n%q\nwant:\n%q\nall-got:\n%s\nall-want:\n%s\n",
					n+1, i+1, format, got[i], want[i], prettyBlocks(got), prettyBlocks(want))
			}
		} else {
			// Match as message
			if got[i] != want[i] {
				t.Fatalf("test %d: fmt.Sprintf(%s, err) at block %d got != want:\n got: %q\nwant: %q", n+1, format, i+1, got[i], want[i])
			}
		}
	}
}

type wrapper struct {
	wrap func(err error) error
	want []string
}

func prettyBlocks(blocks []string) string {
	var out []string

	for _, b := range blocks {
		out = append(out, fmt.Sprintf("%v", b))
	}

	return "   " + strings.Join(out, "\n   ")
}

func testGenericRecursive(t *testing.T, beforeErr error, beforeWant []string, list []wrapper, maxDepth int) {
	if len(beforeWant) == 0 {
		panic("beforeWant must not be empty")
	}
	for _, w := range list {
		if len(w.want) == 0 {
			panic("want must not be empty")
		}

		err := w.wrap(beforeErr)

		// Copy required cause append(beforeWant, ..) modified beforeWant subtly.
		beforeCopy := make([]string, len(beforeWant))
		copy(beforeCopy, beforeWant)

		beforeWant := beforeCopy
		last := len(beforeWant) - 1
		var want []string

		// Merge two stacks behind each other.
		if strings.ContainsAny(beforeWant[last], "\n") && strings.ContainsAny(w.want[0], "\n") {
			want = append(beforeWant[:last], append([]string{beforeWant[last] + "((?s).*)" + w.want[0]}, w.want[1:]...)...)
		} else {
			want = append(beforeWant, w.want...)
		}

		testFormatCompleteCompare(t, maxDepth, err, "%+v", want, false)
		if maxDepth > 0 {
			testGenericRecursive(t, err, want, list, maxDepth-1)
		}
	}
}

func TestFormatCode(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{"%s", `ConfigurationNotValid error`},
		{"%v", `ConfigurationNotValid error`},
		{"%-v", `^service configuration could not be loaded - #3 \[.*mocks_test.go:34 \(.*errors.loadConfig\)\] \(1000\) ConfigurationNotValid error$`},
		{"%+v", `^service configuration could not be loaded - #3 \[.*mocks_test.go:34 \(.*errors.loadConfig\)\] \(1000\) ConfigurationNotValid error; could not decode configuration data - #2 \[.*mocks_test.go:39 \(.*errors.decodeConfig\)\] \(1001\) Data is not valid JSON; could not read configuration file - #1 \[.*mocks_test.go:44 \(.*errors.readConfig\)\] \(1002\) End of input; read: end of input - #0 read: end of input`},
		{"%#-v", `[{\"caller\":\"#3 /home/lk/workspace/golang/src/imooc/mxshop/pkg/errors/mocks_test.go:34 (imooc/mxshop/pkg/errors.loadConfig)\",\"code\":1000,\"error\":\"service configuration could not be loaded\",\"message\":\"ConfigurationNotValid error\"}]`},
		{"%#+v", `[{\"caller\":\"#3 /home/lk/workspace/golang/src/imooc/mxshop/pkg/errors/mocks_test.go:34 (imooc/mxshop/pkg/errors.loadConfig)\",\"code\":1000,\"error\":\"service configuration could not be loaded\",\"message\":\"ConfigurationNotValid error\"},{\"caller\":\"#2 /home/lk/workspace/golang/src/imooc/mxshop/pkg/errors/mocks_test.go:39 (imooc/mxshop/pkg/errors.decodeConfig)\",\"code\":1001,\"error\":\"could not decode configuration data\",\"message\":\"Data is not valid JSON\"},{\"caller\":\"#1 /home/lk/workspace/golang/src/imooc/mxshop/pkg/errors/mocks_test.go:39 (imooc/mxshop/pkg/errors.readConfig)\",\"code\":1002,\"error\":\"could not read configuration file\",\"message\":\"End of input\"},{\"caller\":\"#0\",\"code\":1,\"error\":\"read: end of input\",\"message\":\"read: end of input\"}]`},
	}

	for i, tt := range tests {
		got := fmt.Sprintf(tt.format, loadConfig())
		if !regexp.MustCompile(tt.want).Match([]byte(got)) {
			t.Errorf("test %d: TestFormatCode:\n got %q\n want %q", i+1, got, tt.want)
		}
	}
}
