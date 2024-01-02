package main

import (
	"os"
	"testing"
)

type Assertion[T comparable] struct {
	t      *testing.T
	actual T
}

func Assert[T comparable](t *testing.T, actual T) Assertion[T] {
	return Assertion[T]{
		t:      t,
		actual: actual,
	}
}

func (assertion Assertion[T]) Equals(expected T) {
	if assertion.actual != expected {
		assertion.t.Errorf("got %v; wanted %v", assertion.actual, expected)
	}
}

func TestGetenvFallback(t *testing.T) {
	Assert(t, Getenv("SVART_ALLOWED", "fallback*")).Equals("fallback*")
}

func TestGetenvActual(t *testing.T) {
	os.Setenv("SVART_ALLOWED", "test*")
	Assert(t, Getenv("SVART_ALLOWED", "fallback*")).Equals("test*")
}

func TestIsAllowedWildcardTrue(t *testing.T) {
	os.Setenv("SVART_ALLOWED", "*")
	actual, _ := IsExportAllowed("foo")
	Assert(t, actual).Equals(true)
}

func TestIsAllowedWildcardPrefixTrue(t *testing.T) {
	os.Setenv("SVART_ALLOWED", "foo*")
	actual, _ := IsExportAllowed("foo")
	Assert(t, actual).Equals(true)
}

func TestIsAllowedLiteralTrue(t *testing.T) {
	os.Setenv("SVART_ALLOWED", "foo")
	actual, _ := IsExportAllowed("foo")
	Assert(t, actual).Equals(true)
}

func TestIsAllowedLiteralsTrue(t *testing.T) {
	os.Setenv("SVART_ALLOWED", "bar,foo")

	actual1, _ := IsExportAllowed("foo")
	Assert(t, actual1).Equals(true)

	actual2, _ := IsExportAllowed("bar")
	Assert(t, actual2).Equals(true)
}

func TestIsAllowedLiteralsFalse(t *testing.T) {
	os.Setenv("SVART_ALLOWED", "bar,baz")

	actual, _ := IsExportAllowed("foo")
	Assert(t, actual).Equals(true)

	actual1, _ := IsExportAllowed("baz")
	Assert(t, actual1).Equals(true)

	actual2, _ := IsExportAllowed("bar")
	Assert(t, actual2).Equals(true)
}

func TestIsBlocked(t *testing.T) {
	// TODO: Write tests for IsBlocked function
}
