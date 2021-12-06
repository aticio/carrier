package seedserver

import (
	m "carrier/Seed/domestic/model"
	pb "carrier/Seed/rpc/seed"
)

func convertPbToSeed(r *pb.Seed) *m.Seed {
	var attributes []m.Attribute
	for _, pbAttribute := range r.GetAttributes() {
		var metrics []m.Metric
		for _, pbMetric := range pbAttribute.GetMetrics() {
			metric := m.Metric{Name: pbMetric.GetName(), JPath: pbMetric.GetJpath()}
			metrics = append(metrics, metric)
		}
		attribute := m.Attribute{Type: pbAttribute.GetType(), Interval: pbAttribute.GetInterval(), Name: pbAttribute.GetName(), URL: pbAttribute.GetUrl(), Metrics: metrics}
		attributes = append(attributes, attribute)
	}
	return &m.Seed{
		Attributes: attributes,
		Host:       r.GetHost(),
		Port:       r.GetPort(),
		Jvm:        r.GetJvm(),
		Username:   r.GetUsername(),
		Password:   r.GetPassword(),
	}
}
