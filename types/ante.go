package types

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx Context, msg any) (newCtx Context, err error)

// AnteDecorator wraps the next AnteHandler to perform custom pre-processing.
type AnteDecorator interface {
	AnteHandle(ctx Context, msg any, next AnteHandler) (newCtx Context, err error)
}

// ChainAnteDecorators ChainDecorator chains AnteDecorators together with each AnteDecorator
// wrapping over the decorators further along chain and returns a single AnteHandler.
//
// NOTE: The first element is outermost decorator, while the last element is innermost
// decorator. Decorator ordering is critical since some decorators will expect
// certain checks and updates to be performed (e.g. the Context) before the decorator
// is run. These expectations should be documented clearly in a CONTRACT docline
// in the decorator's godoc.
//
// NOTE: Any application that uses GasMeter to limit transaction processing cost
// MUST set GasMeter with the FIRST AnteDecorator. Failing to do so will cause
// transactions to be processed with an infinite gasmeter and open a DOS attack vector.
// Use `ante.SetUpContextDecorator` or a custom Decorator with similar functionality.
// Returns nil when no AnteDecorator are supplied.
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
