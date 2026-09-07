package ingester

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsBackpressureError(t *testing.T) {
	// Test 1: Backpressure error should be correctly identified
	var bpErr = &ReplayBackpressureError{
		inUse: 100,
		ceiling: 100,
	}
	require.True(t, IsBackpressureError(bpErr), "Should identify backpressure error")

	// Test 2: Generic error should not be identified as backpressure error
	var genericErr = errors.New("generic error")
	require.False(t, IsBackpressureError(genericErr), "Should not identify generic error as backpressure")

	// Test 3: Nil error should not cause panic
	require.False(t, IsBackpressureError(nil), "Should not identify nil error as backpressure")

	// Test 4: Wrapped error should be correctly identified
	var wrappedErr = errors.Join(bpErr, genericErr)
	require.True(t, IsBackpressureError(wrappedErr), "Should identify wrapped backpressure error")

	// Test 5: Non-backpressure wrapped error should not be identified
	var wrappedGenericErr = errors.Join(genericErr, bpErr)
	require.True(t, IsBackpressureError(wrappedGenericErr), "Should identify backpressure in wrapped error")

	// Test 6: Type assertion should work correctly
	var err error = bpErr
	require.IsType(t, &ReplayBackpressureError{}, err, "Should be able to type assert backpressure error")

	// Test 7: Ensure error messages contain expected information
	errMsg := bpErr.Error()
	require.Contains(t, errMsg, "in use", "Error message should contain 'in use'")
	require.Contains(t, errMsg, "ceiling", "Error message should contain 'ceiling'")
	require.Contains(t, errMsg, "cannot recover", "Error message should contain 'cannot recover'")
}
