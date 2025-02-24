package metrics

import (
	"fmt"
	"strings"

	"github.com/kyma-project/control-plane/components/provisioner/internal/model"
	"github.com/kyma-project/control-plane/components/provisioner/internal/persistence/dberrors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=OperationsStatsGetter
type OperationsStatsGetter interface {
	InProgressOperationsCount() (model.OperationsCount, dberrors.Error)
}

type InProgressOperationsCollector struct {
	statsGetter OperationsStatsGetter

	provisioningDesc            *prometheus.Desc
	provisioningNoInstallDesc   *prometheus.Desc
	deprovisioningDesc          *prometheus.Desc
	deprovisioningNoInstallDesc *prometheus.Desc
	upgradeDesc                 *prometheus.Desc

	log logrus.FieldLogger
}

func NewInProgressOperationsCollector(statsGetter OperationsStatsGetter) *InProgressOperationsCollector {
	return &InProgressOperationsCollector{
		statsGetter: statsGetter,

		provisioningDesc: prometheus.NewDesc(
			buildFQName(model.Provision),
			"The number of provisioning operations in progress",
			[]string{},
			nil),
		provisioningNoInstallDesc: prometheus.NewDesc(
			buildFQName(model.ProvisionNoInstall),
			"The number of provisioning without installation operations in progress",
			[]string{},
			nil),
		deprovisioningDesc: prometheus.NewDesc(
			buildFQName(model.Deprovision),
			"The number of deprovisioning operations in progress",
			[]string{},
			nil),
		deprovisioningNoInstallDesc: prometheus.NewDesc(
			buildFQName(model.DeprovisionNoInstall),
			"The number of deprovisioning without uninstallation operations in progress",
			[]string{},
			nil),
		upgradeDesc: prometheus.NewDesc(
			buildFQName(model.Upgrade),
			"The number of upgrade operations in progress",
			[]string{},
			nil),

		log: logrus.WithField("collector", "in-progress-operations"),
	}
}

func (c *InProgressOperationsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.provisioningDesc
	ch <- c.provisioningNoInstallDesc
	ch <- c.deprovisioningDesc
	ch <- c.deprovisioningNoInstallDesc
	ch <- c.upgradeDesc
}

func (c *InProgressOperationsCollector) Collect(ch chan<- prometheus.Metric) {
	inProgressOpsCounts, err := c.statsGetter.InProgressOperationsCount()
	if err != nil {
		c.log.Errorf("failed to get count of operations in progress while collecting metrics: %s", err.Error())

		return
	}

	c.newMeasure(ch,
		c.provisioningDesc,
		inProgressOpsCounts.Count[model.Provision],
	)
	c.newMeasure(ch,
		c.provisioningNoInstallDesc,
		inProgressOpsCounts.Count[model.ProvisionNoInstall],
	)
	c.newMeasure(ch,
		c.deprovisioningDesc,
		inProgressOpsCounts.Count[model.Deprovision],
	)
	c.newMeasure(ch,
		c.deprovisioningNoInstallDesc,
		inProgressOpsCounts.Count[model.DeprovisionNoInstall],
	)
	c.newMeasure(ch,
		c.upgradeDesc,
		inProgressOpsCounts.Count[model.Upgrade],
	)
}

func (c *InProgressOperationsCollector) newMeasure(ch chan<- prometheus.Metric, desc *prometheus.Desc, value int, labelValues ...string) {
	m, err := prometheus.NewConstMetric(
		desc,
		prometheus.GaugeValue,
		float64(value),
		labelValues...)
	if err != nil {
		c.log.Errorf("unable to register metric %s", err.Error())
		return
	}
	ch <- m
}

func buildFQName(operationType model.OperationType) string {
	operation := strings.ToLower(string(operationType))
	name := fmt.Sprintf("in_progress_%s_operations_total", operation)
	return prometheus.BuildFQName(prometheusNamespace, prometheusSubsystem, name)
}
