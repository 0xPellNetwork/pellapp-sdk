package types

import (
	"errors"
	"testing"
	"context"
)

// MockContext implements Context interface for testing
type MockContext struct {
	Context
}

func NewMockContext() Context {
	return Context{baseCtx: context.Background()}
}

// MockDecorator implements AnteDecorator interface for testing
type MockDecorator struct {
	name      string
	shouldErr bool
}

func (d MockDecorator) AnteHandle(ctx Context, msg any, next AnteHandler) (Context, error) {
	if d.shouldErr {
		return ctx, errors.New(d.name + " error")
	}

	// Add decorator name to context to track execution order
	var executionOrder []string
	if val := ctx.Value("executionOrder"); val != nil {
		executionOrder = val.([]string)
	}
	executionOrder = append(executionOrder, d.name)
	ctx = ctx.WithValue("executionOrder", executionOrder)

	return next(ctx, msg)
}

func TestChainAnteDecorators(t *testing.T) {
	// Test empty decorator chain
	t.Run("EmptyDecorators", func(t *testing.T) {
		handler := ChainAnteDecorators()
		if handler != nil {
			t.Error("Expected nil handler for empty decorator chain, but got non-nil")
		}
	})

	// Test single decorator
	t.Run("SingleDecorator", func(t *testing.T) {
		decorator := MockDecorator{name: "decorator1", shouldErr: false}
		handler := ChainAnteDecorators(decorator)

		ctx := NewMockContext()
		newCtx, err := handler(ctx, "test message")

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		executionOrder := newCtx.Value("executionOrder").([]string)
		if len(executionOrder) != 1 || executionOrder[0] != "decorator1" {
			t.Errorf("Expected execution order [decorator1], got: %v", executionOrder)
		}
	})

	// Test multiple decorators
	t.Run("MultipleDecorators", func(t *testing.T) {
		decorator1 := MockDecorator{name: "decorator1", shouldErr: false}
		decorator2 := MockDecorator{name: "decorator2", shouldErr: false}
		decorator3 := MockDecorator{name: "decorator3", shouldErr: false}

		handler := ChainAnteDecorators(decorator1, decorator2, decorator3)

		ctx := NewMockContext()
		newCtx, err := handler(ctx, "test message")

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		executionOrder := newCtx.Value("executionOrder").([]string)
		expected := []string{"decorator1", "decorator2", "decorator3"}

		if len(executionOrder) != len(expected) {
			t.Errorf("Expected %d decorators to execute, but got: %d", len(expected), len(executionOrder))
		}

		for i, name := range expected {
			if executionOrder[i] != name {
				t.Errorf("Expected decorator %d to be %s, got: %s", i, name, executionOrder[i])
			}
		}
	})

	// Test decorator returning error
	t.Run("DecoratorReturnsError", func(t *testing.T) {
		decorator1 := MockDecorator{name: "decorator1", shouldErr: false}
		decorator2 := MockDecorator{name: "decorator2", shouldErr: true}
		decorator3 := MockDecorator{name: "decorator3", shouldErr: false}

		handler := ChainAnteDecorators(decorator1, decorator2, decorator3)

		ctx := NewMockContext()
		_, err := handler(ctx, "test message")

		if err == nil {
			t.Error("Expected error, but got nil")
		}

		expectedErr := "decorator2 error"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message '%s', got: '%s'", expectedErr, err.Error())
		}
	})
}
