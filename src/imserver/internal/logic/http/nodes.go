package http

import (
	"context"

	"baseGo/src/fecho/echo"
)

func (s *Server) nodesWeighted(c echo.Context) error {
	var arg struct {
		Platform string `form:"platform"`
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	res := s.logic.NodesWeighted(context.TODO(), arg.Platform, c.RealIP())
	return c.JSON(200, res)
}

func (s *Server) nodesInstances(c echo.Context) error {
	res := s.logic.NodesInstances(context.TODO())
	return c.JSON(200, res)
}
