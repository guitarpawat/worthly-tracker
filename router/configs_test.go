//go:build test || unit

package router

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/logs"
	"worthly-tracker/model"
)

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

type ConfigSuite struct {
	suite.Suite
	service configService
}

func (s *ConfigSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
}

func (s *ConfigSuite) SetupTest() {
	s.service = configService{}
}

func (s *ConfigSuite) TestGetHeader_200_WithoutPageName() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/configs/header")

	err := s.service.getHeader(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	var headerConfig model.Header
	err = json.Unmarshal(w.Body.Bytes(), &headerConfig)
	s.Require().NoError(err)
	s.Require().Equal(3, len(headerConfig.Links))
	s.Require().False(headerConfig.Links[0].Highlight)
	s.Require().False(headerConfig.Links[1].Highlight)
	s.Require().False(headerConfig.Links[2].Highlight)
}

func (s *ConfigSuite) TestGetHeader_200_WithPageName() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/configs/header/:currentPage")
	c.SetParamNames("currentPage")
	c.SetParamValues("postRecord")

	err := s.service.getHeader(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	var headerConfig model.Header
	err = json.Unmarshal(w.Body.Bytes(), &headerConfig)
	s.Require().NoError(err)
	s.Require().Equal(3, len(headerConfig.Links))
	s.Require().True(headerConfig.Links[0].Highlight)
	s.Require().False(headerConfig.Links[1].Highlight)
	s.Require().False(headerConfig.Links[2].Highlight)
}

func (s *ConfigSuite) TestGetHeader_200_WithUnknownPageName() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/configs/header/:currentPage")
	c.SetParamNames("currentPage")
	c.SetParamValues("test")

	err := s.service.getHeader(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	var headerConfig model.Header
	err = json.Unmarshal(w.Body.Bytes(), &headerConfig)
	s.Require().NoError(err)
	s.Require().Equal(3, len(headerConfig.Links))
	s.Require().False(headerConfig.Links[0].Highlight)
	s.Require().False(headerConfig.Links[1].Highlight)
	s.Require().False(headerConfig.Links[2].Highlight)
}

func (s *ConfigSuite) TestGetHeader_200_WithEmptyPageName() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/configs/header/:currentPage")
	c.SetParamNames("currentPage")
	c.SetParamValues("")

	err := s.service.getHeader(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	var headerConfig model.Header
	err = json.Unmarshal(w.Body.Bytes(), &headerConfig)
	s.Require().NoError(err)
	s.Require().Equal(3, len(headerConfig.Links))
	s.Require().False(headerConfig.Links[0].Highlight)
	s.Require().False(headerConfig.Links[1].Highlight)
	s.Require().False(headerConfig.Links[2].Highlight)
}
