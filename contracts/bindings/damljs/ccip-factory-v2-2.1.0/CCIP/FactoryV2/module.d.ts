// Generated from ../../CCIP/FactoryV2/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e from '@chainlink/chainlink-canton-contracts/ccip-committee-verifier-v2-2.0.0';
import * as pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f from '@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0';
import * as pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40 from '@chainlink/chainlink-canton-contracts/ccip-rate-limiter-v2-2.0.0';
import * as pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 from '@chainlink/chainlink-canton-contracts/mcms-api-1.0.0';
import * as pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb from '@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0';
import * as pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b from '@chainlink/chainlink-canton-contracts/splice-api-token-holding-v1-1.0.0';
import * as pkg7a9aef5d7e6c613c2c57365e732554f049e66f96a6f1d29ca190f5d767580a5a from '@chainlink/chainlink-canton-contracts/ccip-lock-release-token-pool-v2-2.1.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';
import * as pkgb43a9896a3d6753fa23d25f42b6021a20984d807ca73f0dc08be71b582451e1c from '@chainlink/chainlink-canton-contracts/ccip-runtime-v2-2.1.0';
import * as pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 from '@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0';
import * as pkgd76abab3fd465d7ae412818f9caab6b8a0de61128bd7db0464d65987e456a323 from '@chainlink/chainlink-canton-contracts/ccip-core-v2-2.1.0';
import * as pkgf0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688 from '@chainlink/chainlink-canton-contracts/ccip-executor-v2-2.0.0';
import * as pkgf83b5db150d68564a63d9a5760a8a49323424f58ebbea0e089fe79643a8f29c2 from '@chainlink/chainlink-canton-contracts/link-2.0.0';
import * as pkgfad3d7f03cafedfa5060726d0981d9edd1aa91ebd2e120a2fe3b8b7e9fbf49c3 from '@chainlink/chainlink-canton-contracts/ccip-burn-mint-token-pool-v2-2.1.0';

export declare type CCIPFactory = {
  instanceId: string,
  owner: damlTypes.Party,
  mcmsParty: damlTypes.Party,
  usedInstanceIds: damlTypes.Map<string, boolean>,
  deployedContracts: damlTypes.Map<string, damlTypes.ContractId<{}>>,
  perPartyRouterFactoryDeployed: boolean,
}

export declare interface CCIPFactoryInterface {
  Archive: 
    damlTypes.Choice<CCIPFactory, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployBurnMintTokenPool: 
    damlTypes.Choice<CCIPFactory, DeployBurnMintTokenPool, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployCommitteeVerifier: 
    damlTypes.Choice<CCIPFactory, DeployCommitteeVerifier, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployExecutor: 
    damlTypes.Choice<CCIPFactory, DeployExecutor, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployFeeQuoter: 
    damlTypes.Choice<CCIPFactory, DeployFeeQuoter, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployGlobalConfig: 
    damlTypes.Choice<CCIPFactory, DeployGlobalConfig, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployLinkToken: 
    damlTypes.Choice<CCIPFactory, DeployLinkToken, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployLockReleaseTokenPool: 
    damlTypes.Choice<CCIPFactory, DeployLockReleaseTokenPool, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployOffRamp: 
    damlTypes.Choice<CCIPFactory, DeployOffRamp, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployOnRamp: 
    damlTypes.Choice<CCIPFactory, DeployOnRamp, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployPerPartyRouterFactory: 
    damlTypes.Choice<CCIPFactory, DeployPerPartyRouterFactory, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployRMNRemote: 
    damlTypes.Choice<CCIPFactory, DeployRMNRemote, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployRateLimiter: 
    damlTypes.Choice<CCIPFactory, DeployRateLimiter, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  DeployTokenAdminRegistry: 
    damlTypes.Choice<CCIPFactory, DeployTokenAdminRegistry, DeployResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  GetFactoryState: 
    damlTypes.Choice<CCIPFactory, GetFactoryState, FactoryState, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
  SetOwnerToMCMS: 
    damlTypes.Choice<CCIPFactory, SetOwnerToMCMS, damlTypes.ContractId<CCIPFactory>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CCIPFactory, undefined>>;
}
export declare const CCIPFactory:
  damlTypes.Template<CCIPFactory, undefined, '#ccip-factory-v2:CCIP.FactoryV2:CCIPFactory'> &
  damlTypes.ToInterface<CCIPFactory, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  CCIPFactoryInterface

export declare type DeployBurnMintTokenPool = {
  contract: pkgfad3d7f03cafedfa5060726d0981d9edd1aa91ebd2e120a2fe3b8b7e9fbf49c3.CCIP.BurnMintTokenPoolV2.BurnMintTokenPool,
}

export declare const DeployBurnMintTokenPool:
  damlTypes.Serializable<DeployBurnMintTokenPool>

export declare type DeployBurnMintTokenPoolParams = {
  instanceId: string,
  poolOwner: damlTypes.Party,
  ccipOwner: damlTypes.Party,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  decimals: damlTypes.Int,
  rateLimitAdmin: damlTypes.Optional<damlTypes.Party>,
  tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  transferTimeout: pkgfad3d7f03cafedfa5060726d0981d9edd1aa91ebd2e120a2fe3b8b7e9fbf49c3.CCIP.BurnMintTokenPoolV2Types.TransferTimeout,
}

export declare const DeployBurnMintTokenPoolParams:
  damlTypes.Serializable<DeployBurnMintTokenPoolParams>

export declare type DeployCommitteeVerifier = {
  contract: pkg0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e.CCIP.CommitteeVerifierV2.CommitteeVerifier,
}

export declare const DeployCommitteeVerifier:
  damlTypes.Serializable<DeployCommitteeVerifier>

export declare type DeployCommitteeVerifierParams = {
  instanceId: string,
  owner: damlTypes.Party,
  ccipOwner: damlTypes.Party,
  versionTag: string,
  allowListAdmin: damlTypes.Optional<damlTypes.Party>,
  messageSentObservers: damlTypes.Party[],
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  storageLocations: string[],
  storageLocationsAdmin: damlTypes.Party,
  pendingStorageLocationsAdmin: damlTypes.Party,
}

export declare const DeployCommitteeVerifierParams:
  damlTypes.Serializable<DeployCommitteeVerifierParams>

export declare type DeployExecutor = {
  contract: pkgf0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688.CCIP.ExecutorV2.Executor,
}

export declare const DeployExecutor:
  damlTypes.Serializable<DeployExecutor>

export declare type DeployExecutorParams = {
  instanceId: string,
  owner: damlTypes.Party,
  maxCCVsPerMsg: damlTypes.Int,
  allowedFinalityConfig: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig,
  ccvAllowlistEnabled: boolean,
}

export declare const DeployExecutorParams:
  damlTypes.Serializable<DeployExecutorParams>

export declare type DeployFeeQuoter = {
  contract: pkgd76abab3fd465d7ae412818f9caab6b8a0de61128bd7db0464d65987e456a323.CCIP.CoreV2.FeeQuoter.FeeQuoter,
}

export declare const DeployFeeQuoter:
  damlTypes.Serializable<DeployFeeQuoter>

export declare type DeployFeeQuoterParams = {
  instanceId: string,
  linkTokenInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
}

export declare const DeployFeeQuoterParams:
  damlTypes.Serializable<DeployFeeQuoterParams>

export declare type DeployGlobalConfig = {
  contract: pkgd76abab3fd465d7ae412818f9caab6b8a0de61128bd7db0464d65987e456a323.CCIP.CoreV2.GlobalConfig.GlobalConfig,
}

export declare const DeployGlobalConfig:
  damlTypes.Serializable<DeployGlobalConfig>

export declare type DeployGlobalConfigParams = {
  instanceId: string,
  chainSelector: damlTypes.Numeric,
}

export declare const DeployGlobalConfigParams:
  damlTypes.Serializable<DeployGlobalConfigParams>

export declare type DeployLinkToken = {
  contract: pkgf83b5db150d68564a63d9a5760a8a49323424f58ebbea0e089fe79643a8f29c2.Link.Token.LinkRegistry,
}

export declare const DeployLinkToken:
  damlTypes.Serializable<DeployLinkToken>

export declare type DeployLinkTokenParams = {
  instanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
}

export declare const DeployLinkTokenParams:
  damlTypes.Serializable<DeployLinkTokenParams>

export declare type DeployLockReleaseTokenPool = {
  contract: pkg7a9aef5d7e6c613c2c57365e732554f049e66f96a6f1d29ca190f5d767580a5a.CCIP.LockReleaseTokenPoolV2.LockReleaseTokenPool,
}

export declare const DeployLockReleaseTokenPool:
  damlTypes.Serializable<DeployLockReleaseTokenPool>

export declare type DeployLockReleaseTokenPoolParams = {
  instanceId: string,
  poolOwner: damlTypes.Party,
  ccipOwner: damlTypes.Party,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  decimals: damlTypes.Int,
  rateLimitAdmin: damlTypes.Optional<damlTypes.Party>,
  tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  transferTimeout: pkg7a9aef5d7e6c613c2c57365e732554f049e66f96a6f1d29ca190f5d767580a5a.CCIP.LockReleaseTokenPoolV2Types.TransferTimeout,
}

export declare const DeployLockReleaseTokenPoolParams:
  damlTypes.Serializable<DeployLockReleaseTokenPoolParams>

export declare type DeployOffRamp = {
  contract: pkgb43a9896a3d6753fa23d25f42b6021a20984d807ca73f0dc08be71b582451e1c.CCIP.RuntimeV2.OffRamp.OffRamp,
}

export declare const DeployOffRamp:
  damlTypes.Serializable<DeployOffRamp>

export declare type DeployOffRampParams = {
  instanceId: string,
  globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
}

export declare const DeployOffRampParams:
  damlTypes.Serializable<DeployOffRampParams>

export declare type DeployOnRamp = {
  contract: pkgb43a9896a3d6753fa23d25f42b6021a20984d807ca73f0dc08be71b582451e1c.CCIP.RuntimeV2.OnRamp.OnRamp,
}

export declare const DeployOnRamp:
  damlTypes.Serializable<DeployOnRamp>

export declare type DeployOnRampParams = {
  instanceId: string,
  globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  maxUSDCentsPerMsg: damlTypes.Numeric,
}

export declare const DeployOnRampParams:
  damlTypes.Serializable<DeployOnRampParams>

export declare type DeployPerPartyRouterFactory = {
  contract: pkgb43a9896a3d6753fa23d25f42b6021a20984d807ca73f0dc08be71b582451e1c.CCIP.RuntimeV2.PerPartyRouter.PerPartyRouterFactory,
}

export declare const DeployPerPartyRouterFactory:
  damlTypes.Serializable<DeployPerPartyRouterFactory>

export declare type DeployPerPartyRouterFactoryParams = {
  instanceId: string,
  onRamp: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  offRamp: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
}

export declare const DeployPerPartyRouterFactoryParams:
  damlTypes.Serializable<DeployPerPartyRouterFactoryParams>

export declare type DeployRMNRemote = {
  contract: pkgd76abab3fd465d7ae412818f9caab6b8a0de61128bd7db0464d65987e456a323.CCIP.CoreV2.RMNRemote.RMNRemote,
}

export declare const DeployRMNRemote:
  damlTypes.Serializable<DeployRMNRemote>

export declare type DeployRMNRemoteParams = {
  instanceId: string,
  rmnOwner: damlTypes.Party,
  ccipOwner: damlTypes.Party,
  customObservers: damlTypes.Party[],
  cursedSubjects: string[],
}

export declare const DeployRMNRemoteParams:
  damlTypes.Serializable<DeployRMNRemoteParams>

export declare type DeployRateLimiter = {
  contract: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimiter,
}

export declare const DeployRateLimiter:
  damlTypes.Serializable<DeployRateLimiter>

export declare type DeployRateLimiterParams = {
  instanceId: string,
  poolInstanceId: string,
  poolOwner: damlTypes.Party,
  remoteChainSelector: damlTypes.Numeric,
  direction: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimitDirection,
  mode: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimitMode,
  isEnabled: boolean,
  capacity: damlTypes.Numeric,
  rate: damlTypes.Numeric,
}

export declare const DeployRateLimiterParams:
  damlTypes.Serializable<DeployRateLimiterParams>

export declare type DeployResult = {
  factoryCid: damlTypes.ContractId<CCIPFactory>,
  deployedCid: damlTypes.ContractId<{}>,
}

export declare const DeployResult:
  damlTypes.Serializable<DeployResult>

export declare type DeployTokenAdminRegistry = {
  contract: pkgd76abab3fd465d7ae412818f9caab6b8a0de61128bd7db0464d65987e456a323.CCIP.CoreV2.TokenAdminRegistry.TokenAdminRegistry,
}

export declare const DeployTokenAdminRegistry:
  damlTypes.Serializable<DeployTokenAdminRegistry>

export declare type DeployTokenAdminRegistryParams = {
  instanceId: string,
}

export declare const DeployTokenAdminRegistryParams:
  damlTypes.Serializable<DeployTokenAdminRegistryParams>

export declare type FactoryState = {
  instanceId: string,
  owner: damlTypes.Party,
  mcmsParty: damlTypes.Party,
  usedInstanceIds: damlTypes.Map<string, boolean>,
  deployedContracts: damlTypes.Map<string, damlTypes.ContractId<{}>>,
  perPartyRouterFactoryDeployed: boolean,
}

export declare const FactoryState:
  damlTypes.Serializable<FactoryState>

export declare type GetFactoryState = {
  caller: damlTypes.Party,
}

export declare const GetFactoryState:
  damlTypes.Serializable<GetFactoryState>

export declare type SetOwnerToMCMS = {
}

export declare const SetOwnerToMCMS:
  damlTypes.Serializable<SetOwnerToMCMS>
