package helpers

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

var ErrInvocation = errors.New("invocation error")

func RequireArgsLen(method string, args starlark.Tuple, n int) error {
	if args.Len() != n {
		return fmt.Errorf("%w: %s requires %d arg(s), received %d", ErrInvocation, method, n, args.Len())
	}

	return nil
}

func RequireStringArg(method string, args starlark.Tuple, index int, argname string) (string, error) {
	if err := RequireArgsLen(method, args, index+1); err != nil {
		return "", err
	}

	value, ok := starlark.AsString(args.Index(index))
	if ok {
		return value, nil
	}

	return "", fmt.Errorf("%w: %s requires that the '%s' arg be a string", ErrInvocation, method, argname)
}

func OptionalArgs(kwargs []starlark.Tuple) starlark.StringDict {
	dict := make(starlark.StringDict, len(kwargs))

	for _, tuple := range kwargs {
		name, _ := starlark.AsString(tuple.Index(0))
		value := tuple.Index(1)

		dict[name] = value
	}

	return dict
}
