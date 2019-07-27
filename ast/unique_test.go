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
	SetKey(msg, ctx)

	key := GetKey(msg)
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

	SetKey(i32, ctx)
	SetKey(msg, ctx)

	require.NotEqual(t, GetKey(i32), GetKey(msg))
	t.Log(GetKey(msg))
}
