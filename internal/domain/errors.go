package domain

import "errors"

var (
	// ErrUserNotFound es retornado cuando no se encuentra un usuario
	ErrUserNotFound = errors.New("usuario no encontrado")

	// ErrEmailAlreadyExists es retornado cuando se intenta crear un usuario con un email que ya existe
	ErrEmailAlreadyExists = errors.New("el email ya está registrado")

	// ErrInvalidInput es retornado cuando los datos de entrada son inválidos
	ErrInvalidInput = errors.New("datos de entrada inválidos")
)
