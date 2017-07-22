package iso

import (
	"fmt"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepDownload struct {
	step *common.StepDownload
}

func (s *stepDownload) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say(fmt.Sprintf("Wrapper %s", s.step.Url))

	return s.step.Run(state)
}

func (s *stepDownload) Cleanup(multistep.StateBag) {}
