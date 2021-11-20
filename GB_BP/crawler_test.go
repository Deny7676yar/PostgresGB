package main

import (
	mock_internal "github.com/Deny7676yar/Go_level2/GB_BP/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCrawler_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	req := mock_internal.NewMockPage(ctrl)

	gomock.InOrder(
		req.EXPECT().
			GetLinks().
			Return([]string{"GetLinks"}),
		)
}
