package rep

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/onsi/say"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/auction/communication/http/auction_http_client"
	"github.com/pivotal-golang/lager"
)

func RepState(out io.Writer) (err error) {
	client := auction_http_client.New(&http.Client{
		Timeout: 5 * time.Second,
	}, "", "http://localhost:1800", lager.NewLogger("veritas"))

	t := time.Now()
	state, err := client.State()
	dt := time.Since(t)

	if err != nil {
		say.Println(0, "Cell State [%s] - Error:%s", dt, say.Red(err.Error()))
		return err
	}

	name := say.Green("Cell State")
	if state.Evacuating {
		name = say.Red("Cell State - EVAC -")
	}

	rootFSes := []string{}
	for key := range state.RootFSProviders {
		if key != "preloaded" {
			rootFSes = append(rootFSes, say.Yellow(key))
		}
	}

	for key := range state.RootFSProviders["preloaded"].(auctiontypes.FixedSetRootFSProvider).FixedSet {
		rootFSes = append(rootFSes, say.Green("preloaded:%s", key))
	}

	say.Println(0, "%s [%s] - Zone:%s | %s Tasks, %s LRPs | C:%d/%d M:%d/%d D:%d/%d | %s",
		name,
		dt,
		say.Cyan(state.Zone),
		say.Cyan("%d", len(state.Tasks)),
		say.Cyan("%d", len(state.LRPs)),
		state.AvailableResources.Containers,
		state.TotalResources.Containers,
		state.AvailableResources.MemoryMB,
		state.TotalResources.MemoryMB,
		state.AvailableResources.DiskMB,
		state.TotalResources.DiskMB,
		strings.Join(rootFSes, ", "),
	)

	return nil
}
