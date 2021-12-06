package scraperserver

import (
	pb "carrier/Scraper/rpc/scraper"
	pbSeed "carrier/Seed/rpc/seed"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

type Server struct{}

func NewScraperServer() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) GetServerSeed(ctx context.Context, r *pb.Server) (*pbSeed.Seed, error) {
	mBeans, err := getMBeans(r.GetHost(), r.GetPort())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	metricConfs, err := readYAML()
	if err != nil {
		return nil, err
	}
	jvmName, err := getJvmName(mBeans, r.GetHost(), r.GetPort())
	if err != nil {
		return nil, err
	}
	seed := RequestMetrics(metricConfs, mBeans, r.GetHost(), r.GetPort())
	seed.Host = r.GetHost()
	seed.Port = r.GetPort()
	seed.Jvm = jvmName
	seed.Username = viper.GetString("libertyUsername")
	seed.Password = viper.GetString("libertyPassword")
	return seed, nil
}
