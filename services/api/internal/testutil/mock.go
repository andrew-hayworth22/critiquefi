package testutil

import "testing"

// Method is a method that will be mocked
type Method string

// Call represents a mocked call to a method
type Call struct {
	Returns []any
}

// Mock represents a mocked interface
type Mock struct {
	t     *testing.T
	calls map[Method][]Call
}

func NewMock(t *testing.T) Mock {
	t.Helper()
	return Mock{t: t, calls: make(map[Method][]Call)}
}

// On sets the return values for the Next call of a method
func (s *Mock) On(method Method, returns ...any) *Mock {
	s.t.Helper()
	s.calls[method] = append(s.calls[method], Call{Returns: returns})
	return s
}

// Next gets the return values for the Next call of a method
func (s *Mock) Next(method Method) Call {
	s.t.Helper()
	calls, ok := s.calls[method]
	if !ok || len(calls) == 0 {
		s.t.Fatalf("unexpected call to %s", method)
	}
	call := calls[0]
	s.calls[method] = calls[1:]
	return call
}
