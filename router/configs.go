package router

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"slices"
	"worthly-tracker/model"
)

func configsRouter(api *echo.Group) {
	cs := configService{}
	api.GET("/header/:currentPage", cs.getHeader)
	api.GET("/header", cs.getHeader)
}

type configService struct {
}

//	@Summary		Get header configuration
//	@Tags			config
//	@Description	Get header configuration data and determine the link to highlight according to current page
//	@Param			currentPage	path	string	false	"Specified current page"
//	@Produce		json
//	@Success		200	{object}	model.Header	"Success to get header config"
//	@Failure		500	{object}	nil				"Generic server error"
//	@Router			/api/configs/header/{currentPage} [get]
//	@Router			/api/configs/header [get]
func (s *configService) getHeader(c echo.Context) error {
	pageNameParam := c.Param("currentPage")

	var headerConfig = model.Header{
		Title: "Worthly Tracker",
		Links: []model.TopLink{
			{
				Name:       "Records",
				Href:       "/",
				ChildNodes: []model.Link{},
				PageName:   []string{"getRecord", "postRecord"},
			},
			{
				Name: "Reports",
				Href: "#",
				ChildNodes: []model.Link{
					{
						Name: "Report by date",
						Href: "/report_by_date",
					},
					{
						Name: "Report net worth",
						Href: "/report_net_worth",
					},
				},
				PageName: []string{},
			},
			{
				Name: "Settings",
				Href: "#",
				ChildNodes: []model.Link{
					{
						Name: "Manage asset types",
						Href: "/asset_type_setting",
					},
					{
						Name: "Asset type sequence",
						Href: "/asset_type_sequence",
					},
					{
						Name: "Manage assets",
						Href: "/asset_setting",
					},
					{
						Name: "Asset sequence",
						Href: "/asset_sequence",
					},
				},
				PageName: []string{},
			},
		},
	}

	for i, link := range headerConfig.Links {
		if slices.Contains(link.PageName, pageNameParam) {
			headerConfig.Links[i].Highlight = true
			headerConfig.Links[i].Href = "#"
			break
		}
	}

	return c.JSON(http.StatusOK, headerConfig)
}
