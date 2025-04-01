package assert_test

import (
	"fmt"
	"strings"
)

// mockTestingT is a mock implementation of require.TestingT and tHelper
// for testing the require functions themselves.
type mockTestingT struct {
	errorMessages []string
	logMessages   []string
	failed        bool
}

func (m *mockTestingT) Reset() {
	m.errorMessages = nil
	m.logMessages = nil
	m.failed = false
}

func (m *mockTestingT) Log(args ...any) {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, fmt.Sprint(arg))
	}
	m.logMessages = append(m.logMessages, strings.Join(strArgs, " "))
}

func (m *mockTestingT) Logf(format string, args ...any) {
	m.logMessages = append(m.logMessages, fmt.Sprintf(format, args...))
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.errorMessages = append(m.errorMessages, fmt.Sprintf(format, args...))
	m.failed = true
}

func (m *mockTestingT) Helper() {}

func (m *mockTestingT) Failed() bool {
	return m.failed
}

func (m *mockTestingT) ErrorContains(substring string) bool {
	for _, msg := range m.errorMessages {
		if strings.Contains(msg, substring) {
			return true
		}
	}
	return false
}
