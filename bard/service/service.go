package service
import (
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bard/storage"
	"github.com/tamtam-im/logx"
	"fmt"
	"github.com/akaspin/bar/barc/lists"
	"sync"
)

// RPC service
type Service struct {
	Info *proto.ServerInfo
	Storage storage.Storage
}


// Upload spec
func (s *Service) UploadSpec(spec *proto.Spec, res *struct{}) (err error) {

	ok, err := s.Storage.IsSpecExists(spec.ID)
	if err != nil {
		return
	}
	if ok {
		return
	}

	for n, m := range spec.BLOBs {
		exists, err := s.Storage.IsBLOBExists(m)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("blob %s:%s not on bard for spec %s", n, m, spec.ID)
		}
	}

	if err = s.Storage.WriteSpec(*spec); err != nil {
		logx.Error(err)
		return
	}
	logx.Debugf("spec %s stored", spec.ID)
	return
}

// Get all links by spec-id
func (s *Service) GetSpec(id *proto.ID, res *lists.BlobMap) (err error) {

	spec, err := s.Storage.ReadSpec(*id)
	if err != nil {
		return
	}

	res1 := lists.BlobMap{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for name, manifestID := range spec.BLOBs {
		wg.Add(1)
		go func(n string, mID proto.ID) {
			defer wg.Done()

			var err1 error
			mu.Lock()
			rr, err1 := s.Storage.ReadManifest(mID)
			if err1 != nil {
				errs = append(errs, err1)
				return
			}
			res1[n] = *rr
			mu.Unlock()
		}(name, manifestID)
	}
	wg.Wait()

	if len(errs) > 0 {
		err = fmt.Errorf("errors while collecting manifests for spec %s: %s",
			spec.ID, errs)
		return
	}

	*res = res1
	return
}

