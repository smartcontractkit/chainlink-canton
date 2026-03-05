package soliditynode

import (
	"github.com/fbsobreira/gotron-sdk/pkg/address"
)

type GetAccountRequest struct {
	Address string `json:"address"` // account address as a hex string
	Visible bool   `json:"visible"` // Optional,whether the address is in base58 format
}

type Frozen struct {
	FrozenBalance int64 `json:"frozen_balance"` // In Stake 1.0, the total amount of TRX staked by the account to obtain bandwidth
	ExpireTime    int64 `json:"expire_time"`    // In Stake 1.0, the expiration time of the stake operation performed by the account to obtain bandwidth. The account can perform the unstake operation after that time.
}
type FrozenBalanceForEnergy struct {
	FrozenBalance int64 `json:"frozen_balance"` // In Stake 1.0, the total amount of TRX staked by the account to obtain energy
	ExpireTime    int64 `json:"expire_time"`    // In Stake 1.0, the expiration time of the stake operation performed by the account to obtain energy. The account can perform the unstake operation after that time.
}

type AccountResource struct {
	DelegatedFrozenBalanceForEnergy           int64 `json:"delegated_frozen_balance_for_energy"`            // In Stake 1.0, the total amount of TRX staked by the account for others to get energy
	AcquiredDelegatedFrozenBalanceForEnergy   int64 `json:"acquired_delegated_frozen_balance_for_energy"`   // In Stake 1.0, the total amount of TRX staked by other accounts for this account to get energy
	DelegatedFrozenv2BalanceForEnergy         int64 `json:"delegated_frozenV2_balance_for_energy"`          // In Stake 2.0, the total amount of TRX staked by the account for others to get energy
	AcquiredDelegatedFrozenv2BalanceForEnergy int64 `json:"acquired_delegated_frozenV2_balance_for_energy"` // In Stake 2.0, the total amount of TRX staked by other accounts for this account to get energy
	EnergyWindowSize                          int64 `json:"energy_window_size"`                             // The number of block times required for to fully recover. If energy_window_optimized is true, energy_window_size's decimal is 3, else is 0. For example, when energy_window_optimized=true, then energy_window_size=1000 means that it takes one block time for energy to fully recovery; if energy_window_optimized=false, then energy_window_size=1000 means that it takes 1000 blocks time for energy to fully recovery
	EnergyWindowOptimized                     bool  `json:"energy_window_optimized"`                        // Whether to optimize energy recover window
	EnergyUsage                               int64 `json:"energy_usage"`                                   // The amount of energy used by the account
	LatestConsumeTimeForEnergy                int64 `json:"latest_consume_time_for_energy"`                 // The last time the account consumed energy
}

type Key struct {
	Address string `json:"address,omitempty"`
	Weight  int64  `json:"weight,omitempty"`
}

type Permission struct {
	Type           string `json:"type"`
	ID             int32  `json:"id,omitempty"` //Owner id=0, Witness id=1, Active id start by 2
	PermissionName string `json:"permission_name,omitempty"`
	Threshold      int64  `json:"threshold,omitempty"`
	ParentId       int32  `json:"parent_id,omitempty"`
	Operations     string `json:"operations,omitempty"` //1 bit 1 contract
	Keys           []*Key `json:"keys,omitempty"`
}

type Account_Frozen struct {
	Balance    int64 `json:"frozen_balance"` // In Stake 1.0, the total amount of TRX staked by the account to obtain bandwidth
	ExpireTime int64 `json:"expire_time"`    // In Stake 1.0, the expiration time of the stake operation performed by the account to obtain bandwidth. The account can perform the unstake operation after that time.
}

type Account_FreezeV2 struct {
	Type   string `json:"type,omitempty"`
	Amount int64  `json:"amount,omitempty"`
}

type Account_UnFreezeV2 struct {
	Type               string `json:"type,omitempty"`
	UnfreezeAmount     int64  `json:"unfreeze_amount,omitempty"`      // the amount of unstaked TRX
	UnfreezeExpireTime int64  `json:"unfreeze_expire_time,omitempty"` // the start time stamp when the unstaked TRX can be withdrawn, in ms
}

type Vote struct {
	VoteAddress string `json:"vote_address,omitempty"` // the super representative address
	VoteCount   int64  `json:"vote_count,omitempty"`   // the number of votes for this super representative
}

type Asset struct {
	Key   string `json:"key,omitempty"`
	Value int64  `json:"value,omitempty"`
}

type GetAccountResponse struct {
	AccountName                                  string               `json:"account_name"` // The name of the account. The account name can be modified through the wallet/updateaccount interface. The account name can only be changed once.
	Address                                      string               `json:"address"`      // Account address
	CreateTime                                   int64                `json:"create_time"`  // Account creation time, i.e. account activation time on the TRON network
	Balance                                      int64                `json:"balance"`      // TRX balance
	Frozen                                       Account_Frozen       `json:"frozen"`
	DelegatedFrozenBalanceForBandwidth           int64                `json:"delegated_frozen_balance_for_bandwidth"`            // In Stake 1.0, the total amount of TRX staked by the account for others to get bandwidth
	AcquiredDelegatedFrozenBalanceForBandwidth   int64                `json:"acquired_delegated_frozen_balance_for_bandwidth"`   // In Stake 1.0, the total amount of TRX staked by other accounts for this account to get bandwidth
	DelegatedFrozenv2BalanceForBandwidth         int64                `json:"delegated_frozenV2_balance_for_bandwidth"`          // In Stake 2.0, the total amount of TRX staked by the account for others to get bandwidth
	AcquiredDelegatedFrozenv2BalanceForBandwidth int64                `json:"acquired_delegated_frozenV2_balance_for_bandwidth"` // In Stake 2.0, the total amount of TRX staked by other accounts for this account to get bandwidth
	AccountResource                              AccountResource      `json:"account_resource"`
	FrozenV2                                     []Account_FreezeV2   `json:"frozenV2"`                 // []	In Stake 2.0, the total amount of TRX staked to obtain various types of resources does not include the delegated TRX
	UnfrozenV2                                   []Account_UnFreezeV2 `json:"unfrozenV2"`               // []	In Stake 2.0, each unstaking information. One of the unstaking information contains three fields: type: resource type; unfreeze_amount: the amount of unstaked TRX; unfreeze_expire_time: the start time stamp when the unstaked TRX can be withdrawn, in ms.
	NetUsage                                     int64                `json:"net_usage"`                // The amount of bandwidth used by the account
	FreeNetUsage                                 int64                `json:"free_net_usage"`           // The amount of free bandwidth used by the account
	NetWindowSize                                int64                `json:"net_window_size"`          // The number of block times required for bandwidth obtained by stake to fully recover. If net_window_optimized is true, net_window_size's decimal is 3, else is 0. For example, when net_window_optimized=true, then net_window_size=1000 means that it takes one block time for bandwidth to fully recovery; if net_window_optimized=false, then net_window_size=1000 means that it takes 1000 blocks time for bandwidth to fully recovery
	NetWindowOptimized                           bool                 `json:"net_window_optimized"`     // Whether to optimize net recover window
	Votes                                        Vote                 `json:"votes"`                    // The number of votes for each Super Representative
	LatestOprationTime                           int64                `json:"latest_opration_time"`     // The last operation time
	LatestConsumeTime                            int64                `json:"latest_consume_time"`      // The last time the account consumed bandwidth
	LatestConsumeFreeTime                        int64                `json:"latest_consume_free_time"` // The last time the account consumed free bandwidth
	IsWitness                                    bool                 `json:"is_witness"`               // Is Super Representative
	Allowance                                    int64                `json:"allowance"`                // The amount of rewards that can be withdrawn for the account
	LatestWithdrawTime                           int64                `json:"latest_withdraw_time"`     // The last time the account has withdrawn the reward, the super representative or user can only withdraw the reward once within 24 hours
	OwnerPermission                              Permission           `json:"owner_permission"`         // owner permissions
	WitnessPermission                            Permission           `json:"witness_permission"`       // witness permissions
	ActivePermission                             []Permission         `json:"active_permission"`        // active permission
	Asset                                        []Asset              `json:"asset"`                    // <string, int64>	The token id and balance of the TRC10 token in the account
	AssetV2                                      []Asset              `json:"assetV2"`                  // <string, int64>	The token id and balance of the TRC10 token in the account. Note, the V2 version is used after allowing token with same name and the proposal has been activated at present.
	AssetIssuedName                              string               `json:"asset_issued_name"`        // The name of the TRC10 token created by the account
	AssetIssuedId                                string               `json:"asset_issued_ID"`          // TRC10 token ID created by the account
	FreeAssetNetUsage                            []Asset              `json:"free_asset_net_usage"`     // <string, int64>	The amount of free bandwidth consumed by account transferring TRC10 tokens
	FreeAssetNetUsagev2                          []Asset              `json:"free_asset_net_usageV2"`   // <string, int64> The amount of free bandwidth consumed by account transferring TRC10 tokens
}

func (tc *Client) GetAccount(accountAddress address.Address) (*GetAccountResponse, error) {
	getAccountResponse := GetAccountResponse{}
	err := tc.Post("/getaccount", &GetAccountRequest{
		Address: accountAddress.String(),
		Visible: true,
	}, &getAccountResponse)
	if err != nil {
		return nil, err
	}

	return &getAccountResponse, nil
}
