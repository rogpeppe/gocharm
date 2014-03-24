package errgo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gc "launchpad.net/gocheck"
)

func Test(t *testing.T) { gc.TestingT(t) }

type errgoSuite struct{}

var _ = gc.Suite(&errgoSuite{})

func (*errgoSuite) TestErrorStringOneAnnotation(c *gc.C) {
	first := fmt.Errorf("first error")
	err := Annotate(first, "annotation")
	c.Assert(err, gc.ErrorMatches, `annotation: first error`)
}

func (*errgoSuite) TestErrorStringTwoAnnotations(c *gc.C) {
	first := fmt.Errorf("first error")
	err := Annotate(first, "annotation")
	err = Annotate(err, "another")
	c.Assert(err, gc.ErrorMatches, `another, annotation: first error`)
}

func (*errgoSuite) TestErrorStringThreeAnnotations(c *gc.C) {
	first := fmt.Errorf("first error")
	err := Annotate(first, "annotation")
	err = Annotate(err, "another")
	err = Annotate(err, "third")
	c.Assert(err, gc.ErrorMatches, `third, another, annotation: first error`)
}

func (*errgoSuite) TestExampleAnnotations(c *gc.C) {
	for _, test := range []struct {
		err      error
		expected string
	}{
		{
			err:      two(),
			expected: "two: one",
		}, {
			err:      three(),
			expected: "three, two: one",
		}, {
			err:      transtwo(),
			expected: "transtwo: translated (one)",
		}, {
			err:      transthree(),
			expected: "transthree: translated (two: one)",
		}, {
			err:      four(),
			expected: "four, transthree: translated (two: one)",
		},
	} {
		c.Assert(test.err.Error(), gc.Equals, test.expected)
	}
}

func (*errgoSuite) TestAnnotatedErrorCheck(c *gc.C) {
	// Look for a file that we know isn't there.
	dir := c.MkDir()
	_, err := os.Stat(filepath.Join(dir, "not-there"))
	c.Assert(os.IsNotExist(err), gc.Equals, true)
	c.Assert(Check(err, os.IsNotExist), gc.Equals, true)

	err = Annotate(err, "wrap it")
	// Now the error itself isn't a 'IsNotExist'.
	c.Assert(os.IsNotExist(err), gc.Equals, false)
	// However if we use the Check method, it is.
	c.Assert(Check(err, os.IsNotExist), gc.Equals, true)
}

func (*errgoSuite) TestGetErrorStack(c *gc.C) {
	for _, test := range []struct {
		err     error
		matches []string
	}{
		{
			err:     fmt.Errorf("testing"),
			matches: []string{"testing"},
		}, {
			err:     Annotate(fmt.Errorf("testing"), "annotation"),
			matches: []string{"testing"},
		}, {
			err:     one(),
			matches: []string{"one"},
		}, {
			err:     two(),
			matches: []string{"one"},
		}, {
			err:     three(),
			matches: []string{"one"},
		}, {
			err:     transtwo(),
			matches: []string{"translated", "one"},
		}, {
			err:     transthree(),
			matches: []string{"translated", "one"},
		}, {
			err:     four(),
			matches: []string{"translated", "one"},
		},
	} {
		stack := GetErrorStack(test.err)
		c.Assert(stack, gc.HasLen, len(test.matches))
		for i, match := range test.matches {
			c.Assert(stack[i], gc.ErrorMatches, match)
		}
	}
}

func (*errgoSuite) TestDetailedErrorStack(c *gc.C) {
	for _, test := range []struct {
		err      error
		expected string
	}{
		{
			err:      one(),
			expected: "one",
		}, {
			err:      two(),
			expected: "two: one [errgo/test_functions_test.go:16, github.com/errgo/errgo.two]",
		}, {
			err: three(),
			expected: "three [errgo/test_functions_test.go:20, github.com/errgo/errgo.three]\n" +
				"two: one [errgo/test_functions_test.go:16, github.com/errgo/errgo.two]",
		}, {
			err: transtwo(),
			expected: "transtwo: translated [errgo/test_functions_test.go:27, github.com/errgo/errgo.transtwo]\n" +
				"one",
		}, {
			err: transthree(),
			expected: "transthree: translated [errgo/test_functions_test.go:34, github.com/errgo/errgo.transthree]\n" +
				"two: one [errgo/test_functions_test.go:16, github.com/errgo/errgo.two]",
		}, {
			err: four(),
			expected: "four [errgo/test_functions_test.go:38, github.com/errgo/errgo.four]\n" +
				"transthree: translated [errgo/test_functions_test.go:34, github.com/errgo/errgo.transthree]\n" +
				"two: one [errgo/test_functions_test.go:16, github.com/errgo/errgo.two]",
		},
	} {
		stack := DetailedErrorStack(test.err, Default)
		c.Check(stack, gc.DeepEquals, test.expected)
	}
}