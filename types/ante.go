package types

// AnteHandle performs pre-processing operations before the message is handled.
// It receives the current context, the message to be processed, and the next handler in the chain.
// The implementation should perform its logic and then call the next handler in the chain,
// or return an error to stop the processing flow.
//
// Parameters:
//   - ctx: The current context containing request information and state
//   - msg: The message being processed
//   - next: The next handler in the decorator chain
//
// Returns:
//   - newCtx: A potentially modified context to pass to subsequent handlers
//   - err: An error if the message should not be processed further, nil otherwise
type AnteHandler func(ctx Context, msg any) (newCtx Context, err error)

// AnteDecorator wraps the next AnteHandler to perform custom pre-processing.
type AnteDecorator interface {
	AnteHandle(ctx Context, msg any, next AnteHandler) (newCtx Context, err error)
}

// ChainAnteDecorators chains AnteDecorators together and returns a single AnteHandler.
// It creates a decorator chain where each AnteDecorator wraps the decorators further along the chain.
// The resulting AnteHandler will execute decorators in order, with each decorator able to perform
// pre-processing, invoke the next handler, and perform post-processing.
//
// NOTE: The first element is the outermost decorator, while the last element is the innermost.
// Decorator ordering is critical since some decorators may expect certain checks or context
// modifications to be performed before they run. These expectations should be documented clearly
// in a CONTRACT docline in the decorator's godoc.
//
// The chain execution stops if any decorator returns an error. If no decorators are supplied,
// nil is returned.
func ChainAnteDecorators(chain ...AnteDecorator) AnteHandler {
	if len(chain) == 0 {
		return nil
	}

	handlerChain := make([]AnteHandler, len(chain)+1)
	// set the terminal AnteHandler decorator
	handlerChain[len(chain)] = func(ctx Context, msg any) (Context, error) {
		return ctx, nil
	}

	for i := range chain {
		ii := i
		handlerChain[ii] = func(ctx Context, msg any) (Context, error) {
			return chain[ii].AnteHandle(ctx, msg, handlerChain[ii+1])
		}
	}

	return handlerChain[0]
}
