package customerrors

import "fmt"

type InsufficientFundsError struct {
}

func (e *InsufficientFundsError) Error() string {
	return "В вашем счете недостаточно средств"
}

type AccountNotFoundError struct {
	Guid string
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf("Счет (guid: %s) не найден", e.Guid)
}

type TransactionNotFoundError struct {
	Guid string
}

func (e *TransactionNotFoundError) Error() string {
	return fmt.Sprintf("Транзакция (guid: %s) не найдена", e.Guid)
}

type UserNotFoundError struct {
	Guid string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("Пользователь (guid: %s) не найден", e.Guid)
}

type UserNotFoundWithPhoneError struct {
	Phone string
}

func (e *UserNotFoundWithPhoneError) Error() string {
	return fmt.Sprintf("Пользователь с номером телефона %s не найден", e.Phone)
}

type UserAlreadyExistsError struct {
	Phone string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("Пользователь с номером телефона %s уже существует", e.Phone)
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("Внутренняя ошибка сервера: %s", e.Message)
}

type InvalidRequestError struct {
}

func (e *InvalidRequestError) Error() string {
	return "Неверный запрос"
}

type InvalidCredentialsError struct {
}

func (e *InvalidCredentialsError) Error() string {
	return "Неверные учетные данные"
}

type InvalidTokenError struct {
}

func (e *InvalidTokenError) Error() string {
	return "Неверный токен"
}

type TokenExpiredError struct {
}

func (e *TokenExpiredError) Error() string {
	return "Токен просрочен"
}

type TokenNotValidYetError struct {
}

func (e *TokenNotValidYetError) Error() string {
	return "Токен еще не действителен"
}

type InvalidPhoneError struct {
	Phone string
}

func (e *InvalidPhoneError) Error() string {
	return fmt.Sprintf("Неверный номер телефона: %s", e.Phone)
}

type InvalidAmountError struct {
	Amount string
}

func (e *InvalidAmountError) Error() string {
	return fmt.Sprintf("Неверная сумма: %s", e.Amount)
}

type InvalidGuidError struct {
	msg string
}

func (e *InvalidGuidError) Error() string {
	return fmt.Sprintf("Неверный guid: %s", e.msg)
}
