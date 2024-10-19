package errorhandler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
func Request(err error, c *fiber.Ctx, message string) error {
	if err != nil {
		log.Println("an error occured:", err)
	}
	c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"success": false,
		"message": message,
		"data":    nil,
	})
	return err
}
