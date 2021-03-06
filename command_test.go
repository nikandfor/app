package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandRunSimple(t *testing.T) {
	var ok bool
	var buf bytes.Buffer
	stderr = &buf

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Args:        Args{},
		Action:      func(*Command) error { ok = true; return nil },
		HelpText: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{{
			Name:        "sub,s,alias",
			Description: "subcommand",
			Action: func(*Command) error {
				assert.Fail(t, "subcommand called")
				return nil
			},
			Flags: []*Flag{
				NewFlag("subflag", 3, "some sub flag"),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag,f,ff", false, "some flag"),
		},
	}
	require.NotNil(t, c.Args)

	err := c.run([]string{"base", "first", "second", "--flag", "-"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, c.Args, Args{"first", "second", "-"})

	assert.True(t, ok, "expected command not called")

	assert.Equal(t, ``, buf.String())
}

func TestCommandRunSub(t *testing.T) {
	var ok bool
	var buf bytes.Buffer
	stderr = &buf

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Action: func(*Command) error {
			assert.Fail(t, "command called")
			return nil
		},
		HelpText: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{{
			Name:        "sub,s,alias",
			Description: "subcommand",
			Args:        Args{},
			Action:      func(*Command) error { ok = true; return nil },
			Flags: []*Flag{
				NewFlag("subflag", 3, "some sub flag"),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag,f,ff", "empty", "some flag"),
			HelpFlag,
		},
	}
	require.NotNil(t, c.Command("sub").Args)

	err := c.run([]string{"base", "sub", "first", "second", "--flag=value", "-", "--subflag", "4"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, c.Command("sub").Args, Args{"first", "second", "-"})

	assert.True(t, ok, "expected command not called")

	assert.Equal(t, ``, buf.String())
}

func TestCommandRunSub2(t *testing.T) {
	var buf bytes.Buffer
	stderr = &buf

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Action: func(*Command) error {
			assert.Fail(t, "command called")
			return nil
		},
		HelpText: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{{
			Name:        "sub,s,alias",
			Description: "subcommand",
			Args:        Args{},
			Action: func(*Command) error {
				assert.Fail(t, "command called")
				return nil
			},
			Flags: []*Flag{
				NewFlag("subflag", 3, "some sub flag"),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag,f,ff", "empty", "some flag"),
			HelpFlag,
		},
	}
	require.NotNil(t, c.Command("sub").Args)

	err := c.run([]string{"base", "sub", "first", "second", "--flag=value", "-", "--subflag", "4", "--nonexisted"}, nil)
	if err != ErrNoSuchFlag && err.(interface{ Unwrap() error }).Unwrap() != ErrNoSuchFlag {
		assert.Fail(t, "bad error: %v", err)
	}
	assert.Equal(t, c.Command("sub").Args, Args{"first", "second", "-"})

	assert.Equal(t, ``, buf.String())
}
