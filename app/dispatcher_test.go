package main

import (
	"context"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestEchoShouldBeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := &testWriter{}
	errWriter := &testWriter{}

	liner := mocks.NewMockLiner(ctrl)
	liner.EXPECT().Readline().Return("echo Hello World", nil)
	liner.EXPECT().Stdout().Return(writer)
	liner.EXPECT().Stderr().Return(errWriter)

	sut := NewDispatcher(liner)

	err := sut.Execute(context.Background())

	require.NoError(t, err)
	require.Equal(t, "Hello World\n", writer.value)
}

type testWriter struct {
	value string
}

func (t *testWriter) Write(p []byte) (n int, err error) {
	t.value = string(p)
	return len(p), nil
}
