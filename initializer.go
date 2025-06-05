package datatypex

import "context"

type (
	Initializer                   interface{ Init() }
	InitializerWithError          interface{ Init() error }
	InitializerByContext          interface{ Init(context.Context) }
	InitializerByContextWithError interface{ Init(context.Context) error }
)

func InitByContext(ctx context.Context, initializer any) error {
	switch v := initializer.(type) {
	case Initializer:
		v.Init()
		return nil
	case InitializerWithError:
		return v.Init()
	case InitializerByContext:
		v.Init(ctx)
		return nil
	case InitializerByContextWithError:
		return v.Init(ctx)
	default:
		return nil
	}
}

func Init(initializer any) error {
	return InitByContext(context.Background(), initializer)
}
