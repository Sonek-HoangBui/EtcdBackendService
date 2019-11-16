package StringMapKV

import (
	"github.com/OpenStars/backendclients/go/tpoststorageservice/thrift/gen-go/OpenStars/Common/TPostStorageService"
)

type TPostStorageServiceIf interface {
	GetData(idpost int64) (*TPostStorageService.TPostItem, error)
	PutData(idpost int64, data *TPostStorageService.TPostItem) error
	RemoveData(idpost int64) error
}