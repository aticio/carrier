package seedserver

import (
	m "carrier/Seed/domestic/model"
	pb "carrier/Seed/rpc/seed"
)

func convertSeedsToPb(seeds []m.Seed) *pb.Seeds {
	pbseeds := []*pb.Seed{}
	for _, seed := range seeds {
		pbseed := pb.Seed{Host: seed.Host, Port: seed.Port, Jvm: seed.Jvm, Username: seed.Username, Password: seed.Password}
		for _, attribute := range seed.Attributes {
			pbattribute := pb.Seed_Attribute{Type: attribute.Type, Name: attribute.Name, Url: attribute.URL, Interval: attribute.Interval}
			for _, metric := range attribute.Metrics {
				pbmetric := pb.Seed_Attribute_Metric{Name: metric.Name, Jpath: metric.JPath}
				pbattribute.Metrics = append(pbattribute.Metrics, &pbmetric)
			}
			pbseed.Attributes = append(pbseed.Attributes, &pbattribute)
		}
		pbseeds = append(pbseeds, &pbseed)
	}
	return &pb.Seeds{
		Seeds: pbseeds,
	}
}
