package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	// These are needed in order to support the authentication methods
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	nsExample = `
	# view the current namespace from your KUBECONFIG alongside all available namespaces
	kubectl ns

	# switch the namespace to foo if foo selects exactly one namespace, otherwise print a filtered list
	kubectl ns foo`
)

// NsOptions provides information required to update the current context
// on a user's KUBECONFIG
type NsOptions struct {
	configFlags *genericclioptions.ConfigFlags
	rawConfig   api.Config
	args        []string

	userSpecifiedNamespace string
	namespaces             *v1.NamespaceList

	genericclioptions.IOStreams
}

// NewNsOptions provides an instance of NsOptions with default values
func NewNsOptions(streams genericclioptions.IOStreams) *NsOptions {
	return &NsOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// NewNsCmd provides a cobra command wrapping NsOptions
func NewNsCmd(streams genericclioptions.IOStreams) *cobra.Command {
	opt := NewNsOptions(streams)

	cmd := &cobra.Command{
		Use:          "ns [new-namespace]",
		Short:        "Display/Switch current namespace",
		Example:      nsExample,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := opt.Complete(c, args); err != nil {
				return err
			}

			if err := opt.Validate(); err != nil {
				return err
			}

			if err := opt.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

// Complete sets all information required for updating the current namespace
func (o *NsOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get namespaces")
	}
	o.namespaces = namespaces

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *NsOptions) Validate() error {
	if len(o.args) > 1 {
		return fmt.Errorf("either one or no arguments are allowed")
	}

	if len(o.args) > 0 {
		o.userSpecifiedNamespace = o.args[0]
	}

	return nil
}

// Run lists all available namespaces, or updates the current namesapce
// based on a provided namespace.
func (o *NsOptions) Run() error {
	selected := []string{}
	for _, ns := range o.namespaces.Items {
		if ns.GetName() == o.userSpecifiedNamespace {
			selected = []string{ns.GetName()}
			break
		}
		if strings.Contains(ns.GetName(), o.userSpecifiedNamespace) {
			selected = append(selected, ns.GetName())
		}
	}
	switch len(selected) {
	case 0:
		return fmt.Errorf("can't change namespace, \"%s\" does not exist", o.userSpecifiedNamespace)
	case 1:
		return o.changeCurrentNs(selected[0])
	}
	return o.printNamespaces(selected)
}

func (o *NsOptions) changeCurrentNs(newNS string) error {
	if err := o.checkContext(); err != nil {
		return err
	}

	currentNs := o.rawConfig.Contexts[o.rawConfig.CurrentContext].Namespace

	if currentNs != newNS {

		o.rawConfig.Contexts[o.rawConfig.CurrentContext].Namespace = newNS
		if err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(),
			o.rawConfig, true); err != nil {
			return err
		}

		fmt.Fprintf(o.Out, "namespace set to \"%s\"\n", newNS)
	}
	return nil
}

func (o *NsOptions) printNamespaces(namespaces []string) error {
	red := color.New(color.FgRed)

	if err := o.checkContext(); err != nil {
		return err
	}
	currentNS := o.rawConfig.Contexts[o.rawConfig.CurrentContext].Namespace

	for _, ns := range namespaces {
		if ns == currentNS {
			red.Fprintf(o.Out, "%s\n", ns)
		} else {
			fmt.Fprintf(o.Out, "%s\n", ns)
		}
	}

	return nil
}

func (o *NsOptions) checkContext() error {
	currentCtx := o.rawConfig.CurrentContext
	if _, ok := o.rawConfig.Contexts[currentCtx]; !ok {
		return fmt.Errorf("current context %s not found anymore in KUBECONFIG", currentCtx)
	}
	return nil
}
