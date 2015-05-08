package client

import (
	"errors"
)

var (
	// ErrNotFound: thrown if a read operation on a Scalaris ring fails
	// because the key did not exist before.
	ErrNotFound = errors.New("scalaris: key not found")

	// ErrAbort: thrown if a the commit of a write operation on a Scalaris ring fails.
	ErrAbort = errors.New("scalaris: abort")

	// ErrConnection: thrown if an operation on a Scalaris ring fails because
	// a connection does not exist or has been disconnected.
	ErrConnection = errors.New("scalaris: connection error")

	// ErrKeyChanged: thrown if a test_and_set operation on a Scalaris ring
	// fails because the old value did not match the expected value.
	ErrKeyChanged = errors.New("scalaris: key changed")

	// ErrNodeNotFound: thrown if a delete operation on a Scalaris ring fails
	// because no Scalaris node was found.
	ErrNodeNotFound = errors.New("scalaris: node not found")

	// ErrNotAList: thrown if a add_del_on_list operation on a scalaris ring
	// fails because the participating values are not lists.
	ErrNotAList = errors.New("scalaris: not a list")

	// ErrNotANumber: thrown if a add_del_on_list operation on a scalaris ring
	// fails because the participating values are not numbers.
	ErrNotANumber = errors.New("scalaris: not a number")

	// ErrTimeout thrown if a read or write operation on a Scalaris ring
	// fails due to a timeout.
	ErrTimeout = errors.New("scalaris: timeout reached")

	// ErrUnknown: thrown during operations on a Scalaris ring, e.g.
	// if an unknown result has been returned.
	ErrUnknown = errors.New("scalaris: unknown error")
)
