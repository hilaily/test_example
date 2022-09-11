package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

/*
type ITestInterface interface {
	Get(name string) string
}
*/

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockITestInterface(ctrl)
	m.EXPECT().Get("test").DoAndReturn(func(name string) string {
		return name + "_mock"
	}).Times(1)

	r := m.Get("test")
	assert.Equal(t, r, "test_mock")
}
