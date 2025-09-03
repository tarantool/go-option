//go:generate go run github.com/tarantool/go-option/cmd/gentypes -ext-code 1 -package internal/test FullMsgpackExtType
//go:generate go run github.com/tarantool/go-option/cmd/gentypes -ext-code 2 -force -package internal/test HiddenTypeAlias
//go:generate go run github.com/tarantool/go-option/cmd/gentypes -ext-code 3 -imports github.com/google/uuid -package internal/test -marshal-func encodeUUID -unmarshal-func decodeUUID uuid.UUID

package main
