package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("short url already exists") //Ошибка что URL уже существует
	ErrGone          = errors.New("URL is deleted")           //Ошибка что URL уже удален
)
