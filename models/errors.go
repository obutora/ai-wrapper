package models

import (
	"errors"
)

// ErrUnsupportedProvider は、サポートされていないプロバイダが指定された場合に返されるエラーです。
var ErrUnsupportedProvider = errors.New("unsupported provider")

// ErrInvalidAPIKey は、無効なAPIキーが指定された場合に返されるエラーです。
var ErrInvalidAPIKey = errors.New("invalid API key")

// ErrInvalidModel は、無効なモデルが指定された場合に返されるエラーです。
var ErrInvalidModel = errors.New("invalid model")

// ErrEmptyMessages は、メッセージが空の場合に返されるエラーです。
var ErrEmptyMessages = errors.New("empty messages")

// ErrAPIRequest は、APIリクエスト中にエラーが発生した場合に返されるエラーです。
var ErrAPIRequest = errors.New("API request error")
