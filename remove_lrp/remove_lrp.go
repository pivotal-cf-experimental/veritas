package remove_lrp

import (
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"

	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	"github.com/pivotal-golang/lager"
)

func RemoveLRP(cluster []string, guid string) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workerpool.NewWorkerPool(10))
	err := adapter.Connect()
	if err != nil {
		return err
	}

	store := bbs.NewVeritasBBS(adapter, timeprovider.NewTimeProvider(), lager.NewLogger("veritas"))

	return store.RemoveDesiredLRPByProcessGuid(guid)
}
