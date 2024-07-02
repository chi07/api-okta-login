package request

import "github.com/labstack/echo/v4"

type LoginRequest struct {
	OktaToken string `json:"oktaToken" validate:"required"`
}

func (r *LoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	return nil
}
