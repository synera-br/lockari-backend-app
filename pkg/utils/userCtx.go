package utils

import (
	"context"
	"fmt"
)

func GetUserID(ctx context.Context) (*string, error) {

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	userIDFromCtx := ctx.Value("UserID")
	if userIDFromCtx == nil {
		return nil, fmt.Errorf("userID not found in context")
	}

	userID, ok := userIDFromCtx.(string)
	if !ok {
		return nil, fmt.Errorf("userID in context is not a string")
	}

	if userID == "" {
		return nil, fmt.Errorf("userID in context is empty")
	}

	if len(userID) < 1 {
		return nil, fmt.Errorf("userID in context is too short")
	}

	return &userID, nil
}
