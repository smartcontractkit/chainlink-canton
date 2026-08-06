package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
)

var onboardingCmd = &cobra.Command{
	Use:   "onboarding",
	Short: "Registry utility onboarding commands",
}

var (
	waitProviderTimeout string
)

func init() {
	onboardingCmd.AddCommand(checkPackagesCmd)
	onboardingCmd.AddCommand(requestProviderServiceCmd)
	onboardingCmd.AddCommand(waitProviderServiceCmd)
	onboardingCmd.AddCommand(onboardRegistrarCmd)
	onboardingCmd.AddCommand(discoverRegistryFactoriesCmd)
	onboardingCmd.AddCommand(requestCredentialServiceCmd)

	waitProviderServiceCmd.Flags().StringVar(&waitProviderTimeout, "timeout", "15m", "Maximum wait for DA operator acceptance")
}

var checkPackagesCmd = &cobra.Command{
	Use:   "check-packages",
	Short: "Verify required Registry utility DARs are on the participant",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		_, participant, err := rt.connect(cmd.Context(), "provider")
		if err != nil {
			return err
		}
		results, err := registry.CheckPackages(cmd.Context(), participant)
		for _, r := range results {
			status := "missing"
			if r.Found {
				status = "ok"
			}
			fmt.Printf("%s  %s  %s\n", status, r.Name, r.Expected)
		}
		if err != nil {
			return err
		}
		fmt.Println("all required Registry packages present")

		return nil
	},
}

var requestProviderServiceCmd = &cobra.Command{
	Use:   "request-provider-service",
	Short: "Submit ProviderServiceRequest (DA operator must accept)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		client, _, err := rt.connect(cmd.Context(), "provider")
		if err != nil {
			return err
		}
		cid, err := registry.RequestProviderService(cmd.Context(), client, registry.OnboardingParties{
			Operator: rt.Config.Parties.Operator,
			Provider: rt.Config.Parties.Provider,
		})
		if err != nil {
			return err
		}
		rt.State.ProviderServiceRequestCID = cid
		if err := rt.saveState(); err != nil {
			return err
		}
		fmt.Println("ProviderServiceRequest CID:", cid)
		fmt.Println("Next: onboarding wait-provider-service")

		return nil
	},
}

var waitProviderServiceCmd = &cobra.Command{
	Use:   "wait-provider-service",
	Short: "Poll until ProviderService exists after DA acceptance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		timeout, err := time.ParseDuration(waitProviderTimeout)
		if err != nil {
			return fmt.Errorf("invalid --timeout: %w", err)
		}
		client, _, err := rt.connect(cmd.Context(), "provider")
		if err != nil {
			return err
		}
		cid, err := registry.WaitForProviderService(cmd.Context(), client, rt.Config.Parties.Provider, timeout)
		if err != nil {
			return err
		}
		rt.State.ProviderServiceCID = cid
		if err := rt.saveState(); err != nil {
			return err
		}
		fmt.Println("ProviderService CID:", cid)
		fmt.Println("Next: onboarding onboard-registrar")

		return nil
	},
}

var onboardRegistrarCmd = &cobra.Command{
	Use:   "onboard-registrar",
	Short: "Provider accepts registrar service request (creates AllocationFactory + TransferRule)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		providerSvcCID := rt.State.ProviderServiceCID
		if providerSvcCID == "" {
			return fmt.Errorf("provider_service_cid missing in state — run wait-provider-service first")
		}
		client, _, err := rt.connect(cmd.Context(), "provider")
		if err != nil {
			return err
		}
		result, err := registry.OnboardRegistrar(cmd.Context(), client, registry.OnboardingParties{
			Operator:  rt.Config.Parties.Operator,
			Provider:  rt.Config.Parties.Provider,
			Registrar: rt.Config.Parties.Registrar,
		}, providerSvcCID)
		if err != nil {
			return err
		}
		rt.State.ProviderConfigurationCID = result.ProviderConfigurationCID
		rt.State.RegistrarServiceRequestCID = result.RegistrarServiceRequestCID
		rt.State.RegistrarServiceCID = result.RegistrarServiceCID
		rt.State.AllocationFactoryCID = result.AllocationFactoryCID
		rt.State.TransferRuleCID = result.TransferRuleCID
		if err := rt.saveState(); err != nil {
			return err
		}
		fmt.Println("RegistrarService CID:", result.RegistrarServiceCID)
		fmt.Println("AllocationFactory CID:", result.AllocationFactoryCID)
		if result.TransferRuleCID != "" {
			fmt.Println("TransferRule CID:", result.TransferRuleCID)
		}
		fmt.Println("Next: issuer create-instrument <id>")

		return nil
	},
}

var discoverRegistryFactoriesCmd = &cobra.Command{
	Use:   "discover-registry-factories",
	Short: "List Registry service contracts visible to the registrar party",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		client, _, err := rt.connect(cmd.Context(), "registrar")
		if err != nil {
			return err
		}
		disc, err := registry.DiscoverFactories(cmd.Context(), client, rt.Config.Parties.Registrar)
		if err != nil {
			return err
		}
		printDiscovery("ProviderService", disc.ProviderService)
		printDiscovery("RegistrarService", disc.RegistrarService)
		printDiscovery("AllocationFactory", disc.AllocationFactory)
		printDiscovery("TransferRule", disc.TransferRule)
		printDiscovery("InstrumentConfiguration", disc.InstrumentConfiguration)

		return nil
	},
}

var requestCredentialServiceCmd = &cobra.Command{
	Use:   "request-credential-service",
	Short: "Request Credential User Service (not required when using empty credential requirements)",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Fprintln(os.Stderr, "Credential User Service onboarding is not automated in this release.")
		fmt.Fprintln(os.Stderr, "Use empty holderRequirements/issuerRequirements on instruments, or follow DA credential utility docs.")

		return fmt.Errorf("request-credential-service not implemented")
	},
}

func printDiscovery(label string, refs []registry.ContractRef) {
	fmt.Printf("\n%s (%d):\n", label, len(refs))
	for _, ref := range refs {
		fmt.Printf("  %s  %s\n", ref.ContractID, ref.TemplateID)
	}
}
