package main

import "fmt"

type ErrorExpOneLet struct{}
//если пользователь ввёл не один символ
func (e *ErrorExpOneLet) Error() string {
	return "Expected One Leter"
}

type ErrorExistsKey struct{
	Key string
}

//на случай, если ключ такой уже есть
func (e *ErrorExistsKey) Error() string {
	return fmt.Sprintf("Exists Key - '%s'", e.Key)
}

type ErrorAbsentKey struct{
}
//на случай, если ключ такой отсутствкет в БД
func (e *ErrorAbsentKey) Error() string {
	return "Absent Key."
}

