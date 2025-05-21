package messages

// Mensajes de error comunes
const (
	// Errores de validación
	ErrRequiredFullName   = "El nombre completo es requerido"
	ErrRequiredEmail      = "El email es requerido"
	ErrRequiredPassword   = "La contraseña es requerida"
	ErrInvalidEmailFormat = "El formato del email es inválido"
	ErrPasswordTooShort   = "La contraseña debe tener al menos 6 caracteres"
	ErrInvalidInput       = "Datos de entrada inválidos"
	ErrDecodeRequestBody  = "Error al decodificar el cuerpo de la solicitud"

	// Errores de usuario
	ErrUserNotFound       = "Usuario no encontrado"
	ErrEmailAlreadyExists = "El email ya está registrado"
	ErrCreateUser         = "Error al crear el usuario"
	ErrUpdateUser         = "Error al actualizar el usuario"
	ErrDeleteUser         = "Error al eliminar el usuario"
	ErrListUsers          = "Error al obtener la lista de usuarios"
	ErrGetUser            = "Error al obtener el usuario"
	ErrRequiredUserID     = "ID de usuario requerido"
)

// Mensajes de éxito
const (
	// Mensajes de usuario
	SuccessUserCreated = "Usuario creado exitosamente"
	SuccessUserUpdated = "Usuario actualizado exitosamente"
	SuccessUserDeleted = "Usuario eliminado exitosamente"
)
