package config

type BoshConfig struct {
	DirectorUUID string `json:"director_uuid"`
	DirectorHost string `json:"director_host"`
}

type AWSBoshConfig struct {
	BoshConfig
	AwsAccessId         string   `json:"aws_access_id"`
	AwsSecretAcccessKey string   `json:"aws_secret_access_key"`
	Route53ZoneNames    []string `json:"route53_zone_names"`
}
