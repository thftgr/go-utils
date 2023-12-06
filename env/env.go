package env

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

//var AWS, _ = AWS{}.Parse()

type AWS struct {
	INSTANCE_ID         string `env:"INSTANCE_ID,default=null"`
	EC2_INSTANCE_ID_URL string `env:"EC2_INSTANCE_ID_URL,default=http://169.254.169.254/latest/meta-data/instance-id"`
}

func (i AWS) Parse() (env *AWS, err error) {
	env = &AWS{}
	err = envconfig.Process(context.Background(), env)
	return
}

//var OS, _ = OS{}.Parse()

type OS struct {
	HOSTNAME string `env:"HOSTNAME,default="`
}

func (i OS) Parse() (env *OS, err error) {
	env = &OS{}
	err = envconfig.Process(context.Background(), env)
	return
}

//var InfluxDB, _ = InfluxDB{}.Parse()

type InfluxDB struct {
	TOKEN  string `env:"INFLUXDB_TOKEN,default="`
	BUCKET string `env:"INFLUXDB_BUCKET,default="`
	ORG    string `env:"INFLUXDB_ORG,default="`
	URL    string `env:"INFLUXDB_URL,default="`
}

func (i InfluxDB) Parse() (env *InfluxDB, err error) {
	env = &InfluxDB{}
	err = envconfig.Process(context.Background(), env)
	return
}
