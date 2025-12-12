package util

import "github.com/gofiber/fiber/v2"

func Response(c *fiber.Ctx, status int, code int, data interface{}, message ...string) error {
	result := fiber.Map{
		"code": code,
		"data": data,
	}
	if len(message) > 0 && message[0] != "" {
		result["message"] = message[0]
	}
	return c.Status(status).JSON(result)
}

func SuccessResponse(c *fiber.Ctx, data interface{}, message ...string) error {
	return Response(c, 200, 0, data, message...)
}

func ErrorResponse(c *fiber.Ctx, code int, message ...string) error {
	return Response(c, code, -1, nil, message...)
}
