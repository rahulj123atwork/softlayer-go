package data_types

type SoftLayer_Security_Certificate struct {
	Id          int    `json:"id"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

type SoftLayer_Security_Certificate_Parameters struct {
	Parameters []SoftLayer_Security_Certificate_Template `json:"parameters"`
}

type SoftLayer_Security_Certificate_Template struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}
