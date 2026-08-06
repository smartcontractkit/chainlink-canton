package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ccip"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "CCIP TokenAdminRegistry linking and validation",
}

func init() {
	operatorCmd.AddCommand(linkTokenToPoolCmd)
	operatorCmd.AddCommand(validateCmd)
	operatorCmd.AddCommand(acceptAdminRoleCmd)
}

var linkTokenToPoolCmd = &cobra.Command{
	Use:   "link-token-to-pool",
	Short: "Propose administrator, accept admin role, and set BurnMint pool on TAR",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		if rt.State.InstrumentID == "" {
			return fmt.Errorf("instrument_id missing in state — run issuer create-instrument first")
		}
		if rt.Config.CCIP.TokenAdminRegistryCID == "" {
			return fmt.Errorf("ccip.token_admin_registry_cid is required in config")
		}
		if rt.Config.CCIP.CCIPParty == "" {
			return fmt.Errorf("ccip.ccip_party is required in config")
		}
		if rt.Config.CCIP.BurnMintPoolInstanceID == "" {
			return fmt.Errorf("ccip.burn_mint_pool_instance_id is required in config")
		}

		instrumentID := splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(rt.Config.Parties.Registrar),
			Id:    types.TEXT(rt.State.InstrumentID),
		}

		// Connect as CCIP party for TAR lookups and ProposeAdministrator disclosures.
		client, _, err := ledger.ConnectDevnet(cmd.Context(), rt.Config, rt.Config.CCIP.CCIPParty)
		if err != nil {
			return fmt.Errorf("connect as ccip party: %w", err)
		}

		tokenConfigCID, tarCID, err := ccip.RegisterTokenPoolViaClient(cmd.Context(), client, ccip.RegisterTokenPoolClientInput{
			TokenAdminRegistryCID: rt.Config.CCIP.TokenAdminRegistryCID,
			InstrumentId:          instrumentID,
			PoolInstanceID:        rt.Config.CCIP.BurnMintPoolInstanceID,
			CcipParty:             rt.Config.CCIP.CCIPParty,
			PoolOwnerParty:        rt.Config.Parties.Registrar,
		})
		if err != nil {
			return err
		}

		rt.State.TokenConfigCID = tokenConfigCID
		rt.State.TokenAdminRegistryCID = tarCID
		if err := rt.saveState(); err != nil {
			return err
		}

		fmt.Println("TokenConfig CID:", tokenConfigCID)
		fmt.Println("TokenAdminRegistry CID:", tarCID)
		fmt.Println("Next: operator validate")

		return nil
	},
}

var acceptAdminRoleCmd = &cobra.Command{
	Use:   "accept-admin-role",
	Short: "Accept TAR admin role only (when ProposeAdministrator already ran)",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("accept-admin-role is included in link-token-to-pool — run that command instead")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Read-only checks: TAR maps instrument to pool and holdings are visible",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rt, err := loadRuntime(configPath, statePath)
		if err != nil {
			return err
		}
		if rt.State.InstrumentID == "" {
			return fmt.Errorf("instrument_id missing in state")
		}
		if rt.Config.CCIP.BurnMintPoolInstanceID == "" {
			return fmt.Errorf("ccip.burn_mint_pool_instance_id is required in config")
		}
		if rt.Config.CCIP.CCIPParty == "" {
			return fmt.Errorf("ccip.ccip_party is required in config")
		}

		instrumentID := splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(rt.Config.Parties.Registrar),
			Id:    types.TEXT(rt.State.InstrumentID),
		}

		client, _, err := rt.connect(cmd.Context(), "registrar")
		if err != nil {
			return err
		}

		if err := ccip.Validate(cmd.Context(), client, instrumentID, rt.Config.CCIP.CCIPParty, rt.Config.CCIP.BurnMintPoolInstanceID); err != nil {
			return err
		}

		fmt.Println("validate: ok")
		fmt.Printf("  instrument admin=%s id=%s\n", instrumentID.Admin, instrumentID.Id)
		fmt.Printf("  pool instance=%s\n", rt.Config.CCIP.BurnMintPoolInstanceID)

		return nil
	},
}
