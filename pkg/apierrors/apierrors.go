package apierrors

import "errors"

var (
	ErrNoAsset = errors.New("asset not found for current platform")
	ErrReturn  = errors.New("unexpected value returned by API")
)
