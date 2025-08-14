package repository

import "errors"

var (
	//ErrAlreadyExists Ошибка что URL уже существует
	ErrAlreadyExists = errors.New("short url already exists")
	//ErrGone Ошибка что URL уже удален
	ErrGone = errors.New("URL is deleted")
)
