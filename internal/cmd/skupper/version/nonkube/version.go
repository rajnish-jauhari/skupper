package nonkube

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/skupperproject/skupper/internal/cmd/skupper/common"
	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"
	"github.com/skupperproject/skupper/internal/images"
	"github.com/skupperproject/skupper/internal/kube/client"
	"github.com/skupperproject/skupper/internal/utils/configs"
	"github.com/skupperproject/skupper/internal/utils/validator"
	skupperv2alpha1 "github.com/skupperproject/skupper/pkg/generated/client/clientset/versioned/typed/skupper/v2alpha1"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

type CmdVersion struct {
	Client     skupperv2alpha1.SkupperV2alpha1Interface
	KubeClient kubernetes.Interface
	CobraCmd   *cobra.Command
	Flags      *common.CommandVersionFlags
	namespace  string
	output     string
	manifest   configs.ManifestManager
}

func NewCmdVersion() *CmdVersion {

	skupperCmd := CmdVersion{}

	return &skupperCmd
}

func (cmd *CmdVersion) NewClient(cobraCommand *cobra.Command, args []string) {
	cli, err := client.NewClient(cobraCommand.Flag("namespace").Value.String(), cobraCommand.Flag("context").Value.String(), cobraCommand.Flag("kubeconfig").Value.String())

	if err == nil {
		cmd.namespace = cli.Namespace
	}
}

func (cmd *CmdVersion) ValidateInput(args []string) error {
	var validationErrors []error
	outputTypeValidator := validator.NewOptionValidator(common.OutputTypes)

	if cmd.Flags != nil && cmd.Flags.Output != "" {
		ok, err := outputTypeValidator.Evaluate(cmd.Flags.Output)
		if !ok {
			validationErrors = append(validationErrors, fmt.Errorf("output type is not valid: %s", err))
		} else {
			cmd.output = cmd.Flags.Output
		}
	}

	return errors.Join(validationErrors...)
}

func (cmd *CmdVersion) InputToOptions() {
	if cmd.output != "" {
		cmd.manifest = configs.ManifestManager{Components: images.NonKubeComponents, EnableSHA: true}
	} else {
		cmd.manifest = configs.ManifestManager{Components: images.NonKubeComponents, EnableSHA: false}
	}

}

func (cmd *CmdVersion) Run() error {
	files := cmd.manifest.GetConfiguredManifest()
	if cmd.output != "" {
		encodedOutput, err := utils.Encode(cmd.output, files)
		if err != nil {
			return err
		}
		fmt.Println(encodedOutput)
	} else {
		tw := tabwriter.NewWriter(os.Stdout, 8, 8, 1, '\t', tabwriter.TabIndent)
		_, _ = fmt.Fprintln(tw, fmt.Sprintf("%s\t%s", "COMPONENT", "VERSION"))

		for _, file := range files.Components {
			fmt.Fprintln(tw, fmt.Sprintf("%s\t%s", file.Component, file.Version))
		}
		_ = tw.Flush()
	}
	return nil
}

func (cmd *CmdVersion) WaitUntil() error { return nil }
