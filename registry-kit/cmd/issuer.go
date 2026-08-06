package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
)

var issuerCmd = &cobra.Command{
	Use:   "issuer",
	Short: "Registry token lifecycle commands",
}

var (
	mintAmount string
	mintHolder string
)

func init() {
	issuerCmd.AddCommand(createInstrumentCmd)
	issuerCmd.AddCommand(mintCmd)
	issuerCmd.AddCommand(acceptMintCmd)
	issuerCmd.AddCommand(querySupplyCmd)
	issuerCmd.AddCommand(issueCredentialCmd)

	mintCmd.Flags().StringVar(&mintAmount, "amount", "", "Amount to mint (required)")
	mintCmd.Flags().StringVar(&mintHolder, "holder", "", "Holder party (default: registrar)")
	_ = mintCmd.MarkFlagRequired("amount")
}

var createInstrumentCmd = &cobra.Command{
	Use:   "create-instrument [instrument-id]",
	Short: "Create InstrumentConfiguration with empty credential requirements",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instrumentID := args[0]
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		registrarSvcCID := rt.State.RegistrarServiceCID
		if registrarSvcCID == "" {
			return fmt.Errorf("registrar_service_cid missing in state — run onboarding onboard-registrar first")
		}
		client, _, err := rt.connect(cmd.Context(), "registrar")
		if err != nil {
			return err
		}
		instCfgCID, err := registry.CreateInstrumentConfiguration(cmd.Context(), client, rt.Config.Parties.Registrar, registrarSvcCID, instrumentID)
		if err != nil {
			return err
		}
		rt.State.InstrumentID = instrumentID
		rt.State.InstrumentConfigurationCID = instCfgCID
		if err := rt.saveState(); err != nil {
			return err
		}
		fmt.Println("InstrumentConfiguration CID:", instCfgCID)
		fmt.Printf("Next: issuer mint --amount <n>  (instrument %s)\n", instrumentID)

		return nil
	},
}

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Request mint via AllocationFactory (DA operator backend context)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		if rt.State.InstrumentID == "" {
			return fmt.Errorf("instrument_id missing in state — run create-instrument first")
		}
		allocCID := rt.State.AllocationFactoryCID
		if allocCID == "" {
			return fmt.Errorf("allocation_factory_cid missing in state — run onboarding onboard-registrar or discover-registry-factories")
		}
		client, _, err := rt.connect(cmd.Context(), "registrar")
		if err != nil {
			return err
		}
		holder := mintHolder
		if holder == "" {
			holder = rt.Config.Parties.Registrar
		}
		mintReqCID, err := registry.RequestMintViaOperatorBackend(cmd.Context(), client, rt.operatorClient(), registry.MintDevnetInput{
			RegistrarParty:       rt.Config.Parties.Registrar,
			InstrumentID:         rt.State.InstrumentID,
			AllocationFactoryCID: allocCID,
			Holder:               holder,
			Amount:               mintAmount,
		})
		if err != nil {
			return err
		}
		rt.State.LastMintRequestCID = mintReqCID
		if err := rt.saveState(); err != nil {
			return err
		}
		fmt.Println("MintRequest CID:", mintReqCID)
		fmt.Println("Next: issuer accept-mint")

		return nil
	},
}

var acceptMintCmd = &cobra.Command{
	Use:   "accept-mint",
	Short: "Accept pending mint request (DA operator backend context)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		mintReqCID := rt.State.LastMintRequestCID
		if mintReqCID == "" {
			client, _, cerr := rt.connect(cmd.Context(), "registrar")
			if cerr != nil {
				return cerr
			}
			mintReqCID, err = registry.FindMintRequestForInstrument(cmd.Context(), client, rt.Config.Parties.Registrar, rt.State.InstrumentID)
			if err != nil {
				return err
			}
			if mintReqCID == "" {
				return fmt.Errorf("no MintRequest CID in state or ACS — run issuer mint first")
			}
		}
		client, _, err := rt.connect(cmd.Context(), "registrar")
		if err != nil {
			return err
		}
		holdingCID, err := registry.AcceptMintViaOperatorBackend(cmd.Context(), client, rt.operatorClient(), rt.Config.Parties.Registrar, mintReqCID)
		if err != nil {
			return err
		}
		fmt.Println("Registry Holding CID:", holdingCID)

		return nil
	},
}

var querySupplyCmd = &cobra.Command{
	Use:   "query-supply",
	Short: "Aggregate Registry Holding balances for the configured instrument",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		if rt.State.InstrumentID == "" {
			return fmt.Errorf("instrument_id missing in state — run create-instrument first")
		}
		client, _, err := rt.connect(cmd.Context(), "registrar")
		if err != nil {
			return err
		}
		total, rows, err := registry.QuerySupply(cmd.Context(), client, rt.Config.Parties.Registrar, rt.State.InstrumentID)
		if err != nil {
			return err
		}
		fmt.Printf("Instrument %s total supply: %s\n", rt.State.InstrumentID, total.String())
		for _, row := range rows {
			fmt.Printf("  owner=%s amount=%s cid=%s\n", row.Owner, row.Amount.String(), row.ContractID)
		}
		if total.IsZero() && len(rows) == 0 {
			fmt.Println("(no holdings found)")
		}

		return nil
	},
}

var issueCredentialCmd = &cobra.Command{
	Use:   "issue-credential",
	Short: "Issue registrar credential (only when instrument requires credentials)",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("issue-credential not implemented — use empty credential requirements or DA credential utility docs")
	},
}
