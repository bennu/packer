package iso

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	vmwcommon "github.com/hashicorp/packer/builder/vmware/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepDownload struct {
	step *common.StepDownload
}

func (s *stepDownload) Run(state multistep.StateBag) multistep.StepAction {
	cache := state.Get("cache").(packer.Cache)
	driver := state.Get("driver").(vmwcommon.Driver)
	ui := state.Get("ui").(packer.Ui)

	esx5, ok := driver.(*ESX5Driver)
	if !ok {
		return multistep.ActionContinue
	}

	ui.Say("Verifying remote cache")

	for _, url := range s.step.Url {
		targetPath := s.step.TargetPath

		if targetPath == "" {
			hash := sha1.Sum([]byte(url))
			cacheKey := fmt.Sprintf("%s.%s", hex.EncodeToString(hash[:]), s.step.Extension)
			targetPath = cache.Lock(cacheKey)
			cache.Unlock(cacheKey)
		}

		remotePath := esx5.cachePath(targetPath)

		if esx5.verifyChecksum(s.step.ChecksumType, s.step.Checksum, remotePath) {
			state.Put(s.step.ResultKey, "skip_upload:"+remotePath)
			ui.Message("Remote cache verified, skipping download step")
			return multistep.ActionContinue
		}
	}

	return s.step.Run(state)
}

func (s *stepDownload) Cleanup(multistep.StateBag) {}
