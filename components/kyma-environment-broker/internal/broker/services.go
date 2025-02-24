package broker

import (
	"context"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/euaccess"

	"fmt"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/middleware"

	"github.com/pivotal-cf/brokerapi/v8/domain"
	"github.com/sirupsen/logrus"
)

const (
	ControlsOrderKey = "_controlsOrder"
	PropertiesKey    = "properties"
)

type ServicesEndpoint struct {
	log            logrus.FieldLogger
	cfg            Config
	servicesConfig ServicesConfig

	enabledPlanIDs map[string]struct{}
}

func NewServices(cfg Config, servicesConfig ServicesConfig, log logrus.FieldLogger) *ServicesEndpoint {
	enabledPlanIDs := map[string]struct{}{}
	for _, planName := range cfg.EnablePlans {
		id := PlanIDsMapping[planName]
		enabledPlanIDs[id] = struct{}{}
	}

	return &ServicesEndpoint{
		log:            log.WithField("service", "ServicesEndpoint"),
		cfg:            cfg,
		servicesConfig: servicesConfig,
		enabledPlanIDs: enabledPlanIDs,
	}
}

// Services gets the catalog of services offered by the service broker
//
//	GET /v2/catalog
func (b *ServicesEndpoint) Services(ctx context.Context) ([]domain.Service, error) {
	var availableServicePlans []domain.ServicePlan
	// we scope to the kymaruntime service only
	class, ok := b.servicesConfig[KymaServiceName]
	if !ok {
		return nil, fmt.Errorf("while getting %s class data", KymaServiceName)
	}

	provider, ok := middleware.ProviderFromContext(ctx)
	platformRegion, ok := middleware.RegionFromContext(ctx)
	for _, plan := range Plans(class.Plans, provider, b.cfg.IncludeAdditionalParamsInSchema, euaccess.IsEURestrictedAccess(platformRegion)) {
		// filter out not enabled plans
		if _, exists := b.enabledPlanIDs[plan.ID]; !exists {
			continue
		}
		// p := plan.PlanDefinition

		availableServicePlans = append(availableServicePlans, plan)
	}

	return []domain.Service{
		{
			ID:                   KymaServiceID,
			Name:                 KymaServiceName,
			Description:          class.Description,
			Bindable:             false,
			InstancesRetrievable: true,
			Tags: []string{
				"SAP",
				"Kyma",
			},
			Plans: availableServicePlans,
			Metadata: &domain.ServiceMetadata{
				DisplayName:         class.Metadata.DisplayName,
				ImageUrl:            class.Metadata.ImageUrl,
				LongDescription:     class.Metadata.LongDescription,
				ProviderDisplayName: class.Metadata.ProviderDisplayName,
				DocumentationUrl:    class.Metadata.DocumentationUrl,
				SupportUrl:          class.Metadata.SupportUrl,
			},
			AllowContextUpdates: true,
		},
	}, nil
}
