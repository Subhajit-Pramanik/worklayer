package tenant

import "github.com/gofiber/fiber/v2"

func getOrgIDFromLocals(c *fiber.Ctx) string {
	orgID, _ := c.Locals("organization_id").(string)
	return orgID
}

func getProjectIDFromLocals(c *fiber.Ctx) string {
	projectID, _ := c.Locals("project_id").(string)
	return projectID
}
