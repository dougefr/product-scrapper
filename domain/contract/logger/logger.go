package logger

import "context"

type (
	// Logger interface do logger da aplicação
	Logger interface {
		Debug(ctx context.Context, message string, valueMaps ...Body)
		Info(ctx context.Context, message string, valueMaps ...Body)
		Warn(ctx context.Context, message string, valueMaps ...Body)
		Error(ctx context.Context, message string, valueMaps ...Body)
		Fatal(ctx context.Context, message string, valueMaps ...Body)
	}

	// Body informações adicionais que serão impressas no JSON
	Body = map[string]interface{}
)
