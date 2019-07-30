package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFieldKey(t *testing.T) {
	msg := &Message{
		File:      nil,
		ParentMsg: nil,
		Name:      "message",
		Fields:    nil,
	}

	ctx := UniqueContext{}
	SetUnique(msg, ctx)

	key := GetUnique(msg)
	require.Equal(t, key+"::Name", GetFieldKey(msg, &msg.Name))

	require.Panics(t, func() {
		GetFieldKey(msg, nil)
	})
	require.Panics(t, func() {
		str := "string"
		GetFieldKey(msg, &str)
	})
	require.Panics(t, func() {
		GetFieldKey(msg, 0)
	})
}

func TestGetKey(t *testing.T) {
	i32 := &Int32{}
	msg := &Message{
		Name: "message",
	}
	ctx := UniqueContext{}

	SetUnique(i32, ctx)
	SetUnique(msg, ctx)

	require.NotEqual(t, GetUnique(i32), GetUnique(msg))
	t.Log(GetUnique(msg))
}
