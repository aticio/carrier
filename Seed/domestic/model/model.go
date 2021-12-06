package model

//Seed struct
type Seed struct {
	Host       string      `bson:"host"`
	Port       string      `bson:"port"`
	Jvm        string      `bson:"jvm"`
	Username   string      `bson:"username"`
	Password   string      `bson:"password"`
	Attributes []Attribute `bson:"attributes"`
}

//Attribute struct
type Attribute struct {
	Type     string   `bson:"type"`
	Name     string   `bson:"name"`
	URL      string   `bson:"url"`
	Interval string   `bson:"interval"`
	Metrics  []Metric `bson:"metrics"`
}

//Metric struct
type Metric struct {
	Name  string `bson:"name"`
	JPath string `bson:"jpath"`
}
