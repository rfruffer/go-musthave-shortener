package repository

import "errors"

// Переменные отвечающие за ошибки
var (
	// ErrAlreadyExists — ошибка, означающая, что короткий URL уже существует.
	ErrAlreadyExists = errors.New("short url already exists")

	// ErrGone — ошибка, означающая, что URL был удалён.
	ErrGone = errors.New("URL is deleted")
)
