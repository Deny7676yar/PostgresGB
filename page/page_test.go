package page

import (
	mock_page "github.com/Deny7676yar/Go_level2/GB_BP/page/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PageStub struct {
	gettitleresult func() string
}

func (ps *PageStub) GetTitleResult() string {
	return ps.gettitleresult()
}

func GetTitle(p *PageStub) string {
	return p.GetTitleResult() + "title"
}

func TestPage_GetTitle1(t *testing.T) {
	stub := &PageStub{func() string {
		return "title"
	}}
	assert.Equal(t, "title", GetTitle(stub))

}

func TestPage_GetTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()



	page := mock_page.NewMockPage(ctrl)

	gomock.InOrder(
		page.EXPECT().
			GetTitle().
			Return("yandex"),
		)

	collector := NewCollector(page)
	getlist, err := collector.
}