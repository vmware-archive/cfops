package storageinterface

type StorageEngine interface {
	GetFeatureStatusValue(featureSignature string) (status string, err error)
	Close() (err error)
}
