package console

import (
	"github.com/gofiber/fiber/v2"
	consoledto "github.com/vyolayer/vyolayer/internal/shared/dto/console"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
)

// @Summary List Project Services
// @Description List all services in a project
// @Tags Project Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectID path string true "Project ID"
// @Success 200 {object} response.SuccessResponse{data=[]consoledto.ServiceManifestDTO}
// @Router /console/projects/{projectID}/services [get]
func (h *ProjectServiceHandler) list(c *fiber.Ctx) error {
	req := &consolev1.ListProjectServicesRequest{
		ProjectId: getProjectIDFromLocals(c),
	}

	resp, err := h.client.ListProjectServices(c.UserContext(), req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	servicesDTO := make([]consoledto.ServiceManifestDTO, len(resp.GetData()))
	for i, service := range resp.GetData() {
		servicesDTO[i] = consoledto.ServiceManifestDTO{
			Key:         service.GetKey(),
			Name:        service.GetName(),
			Status:      service.GetStatus(),
			Plan:        service.GetPlan(),
			Icon:        service.GetIcon(),
			Description: service.GetDescription(),
		}
	}

	return response.Success(c, servicesDTO)
}

// @Summary Get Project Service
// @Description Get details of a project service by its key
// @Tags Project Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectID path string true "Project ID"
// @Param serviceKey path string true "Service Key"
// @Success 200 {object} response.SuccessResponse{data=consoledto.ServiceManifestWithResourcesDTO}
// @Router /console/projects/{projectID}/services/{serviceKey} [get]
func (h *ProjectServiceHandler) get(c *fiber.Ctx) error {
	req := &consolev1.GetProjectServiceManifestRequest{
		ProjectId:  getProjectIDFromLocals(c),
		ServiceKey: c.Params(ParamServiceKey),
	}

	grpcResp, err := h.client.GetProjectServiceManifest(c.UserContext(), req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	grpcData := grpcResp.GetData()
	if grpcData == nil {
		return response.Error(c, errors.NotFound("service not found"))
	}

	var resourcesDTO []consoledto.ResourceDTO
	for _, res := range grpcData.GetResources() {
		var columnsDTO []consoledto.ColumnDTO
		for _, col := range res.GetColumns() {
			columnsDTO = append(columnsDTO, consoledto.ColumnDTO{
				Key:      col.GetKey(),
				Label:    col.GetLabel(),
				Type:     col.GetType(),
				Sortable: col.GetSortable(),
				Visible:  col.GetVisible(),
			})
		}

		var actionsDTO []consoledto.ActionDTO
		for _, act := range res.GetActions() {
			actionsDTO = append(actionsDTO, consoledto.ActionDTO{
				Key:     act.GetKey(),
				Label:   act.GetLabel(),
				Scope:   act.GetScope(),
				Variant: act.GetVariant(),
				Danger:  act.GetDanger(),
			})
		}

		var filtersDTO []consoledto.FilterDTO
		for _, fil := range res.GetFilters() {
			filtersDTO = append(filtersDTO, consoledto.FilterDTO{
				Key:   fil.GetKey(),
				Label: fil.GetLabel(),
				Type:  fil.GetType(),
			})
		}

		resourcesDTO = append(resourcesDTO, consoledto.ResourceDTO{
			Key:     res.GetKey(),
			Label:   res.GetLabel(),
			Route:   res.GetRoute(),
			Icon:    res.GetIcon(),
			Columns: columnsDTO,
			Actions: actionsDTO,
			Filters: filtersDTO,
		})
	}

	serviceManifestDTO := &consoledto.ServiceManifestWithResourcesDTO{
		ServiceManifestDTO: consoledto.ServiceManifestDTO{
			Key:         grpcData.GetKey(),
			Name:        grpcData.GetName(),
			Description: grpcData.GetDescription(),
			Status:      grpcData.GetStatus(),
			Plan:        grpcData.GetPlan(),
			Icon:        grpcData.GetIcon(),
		},
		Resources: resourcesDTO,
	}

	return response.Success(c, serviceManifestDTO)
}
