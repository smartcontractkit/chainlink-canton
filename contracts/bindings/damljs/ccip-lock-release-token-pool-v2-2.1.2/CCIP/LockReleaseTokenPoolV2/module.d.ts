// Generated from ../../CCIP/LockReleaseTokenPoolV2/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg1377c6dda2856b806eba169ee7a1b545da7ad8a9416b53f90fdb389853c9f8e4 from '@chainlink/chainlink-canton-contracts/ccip-rate-limiter-v2-2.0.2';
import * as pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4 from '@chainlink/chainlink-canton-contracts/ccip-extension-api-v2-2.0.0';
import * as pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f from '@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0';
import * as pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 from '@chainlink/chainlink-canton-contracts/mcms-api-1.0.0';
import * as pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb from '@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0';
import * as pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b from '@chainlink/chainlink-canton-contracts/splice-api-token-holding-v1-1.0.0';
import * as pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58 from '@chainlink/chainlink-canton-contracts/ccip-api-v2-2.0.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';
import * as pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 from '@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0';

import * as CCIP_LockReleaseTokenPoolV2Types from '../../CCIP/LockReleaseTokenPoolV2Types/module';

export declare type AddPoolReceiveContextContractValue = {
  contextKey: string,
  referredContract: damlTypes.ContractId<{}>,
}

export declare const AddPoolReceiveContextContractValue:
  damlTypes.Serializable<AddPoolReceiveContextContractValue>

export declare type AddPoolReceiveContextNonContractValue = {
  contextKey: string,
  value: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.AnyValue,
}

export declare const AddPoolReceiveContextNonContractValue:
  damlTypes.Serializable<AddPoolReceiveContextNonContractValue>

export declare type ApplyChainUpdates = {
  remoteChainSelectorsToRemove: damlTypes.Numeric[],
  chainsToAdd: CCIP_LockReleaseTokenPoolV2Types.ChainUpdate[],
}

export declare const ApplyChainUpdates:
  damlTypes.Serializable<ApplyChainUpdates>

export declare type ApplyTokenTransferFeeConfigUpdates = {
  tokenTransferFeeConfigArgs: CCIP_LockReleaseTokenPoolV2Types.TokenTransferFeeConfigArgs[],
  disableTokenTransferFeeConfigArgs: damlTypes.Numeric[],
}

export declare const ApplyTokenTransferFeeConfigUpdates:
  damlTypes.Serializable<ApplyTokenTransferFeeConfigUpdates>

export declare type CalculateFee = {
  tokenAdminRegistryCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenAdminRegistry>,
  tokenConfigCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenConfig>,
  sendingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>,
  feeQuoterCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.FeeQuoter.IFeeQuoter>,
  tokenInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const CalculateFee:
  damlTypes.Serializable<CalculateFee>

export declare type ClearPoolReceiveContext = {
}

export declare const ClearPoolReceiveContext:
  damlTypes.Serializable<ClearPoolReceiveContext>

export declare type GetFee = {
  feeQuoterCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.FeeQuoter.IFeeQuoter>,
  destChainSelector: damlTypes.Numeric,
  tokenInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const GetFee:
  damlTypes.Serializable<GetFee>

export declare type GetRequiredCCVs = {
  remoteChainSelector: damlTypes.Numeric,
  sourceAmount: string,
  finality: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig,
  extraData: string,
  direction: pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.TokenPool.TransferDirection,
  caller: damlTypes.Party,
}

export declare const GetRequiredCCVs:
  damlTypes.Serializable<GetRequiredCCVs>

export declare type LockOrBurn = {
  tokenAdminRegistryCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenAdminRegistry>,
  tokenConfigCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenConfig>,
  rmnRemoteCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote>,
  sendingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>,
  senderInputCids: damlTypes.ContractId<pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding>[],
  amount: damlTypes.Numeric,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const LockOrBurn:
  damlTypes.Serializable<LockOrBurn>

export declare type LockReleaseTokenPool = {
  instanceId: string,
  poolOwner: damlTypes.Party,
  ccipOwner: damlTypes.Party,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  decimals: damlTypes.Int,
  rateLimitAdmin: damlTypes.Optional<damlTypes.Party>,
  remoteChainConfigs: damlTypes.Map<damlTypes.Numeric, CCIP_LockReleaseTokenPoolV2Types.RemoteChainConfig>,
  tokenTransferFeeConfigs: damlTypes.Map<damlTypes.Numeric, CCIP_LockReleaseTokenPoolV2Types.TokenTransferFeeConfig>,
  poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  transferTimeout: CCIP_LockReleaseTokenPoolV2Types.TransferTimeout,
  deps: LockReleaseTokenPoolDeps,
}

export declare interface LockReleaseTokenPoolInterface {
  AddPoolReceiveContextContractValue: 
    damlTypes.Choice<LockReleaseTokenPool, AddPoolReceiveContextContractValue, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  AddPoolReceiveContextNonContractValue: 
    damlTypes.Choice<LockReleaseTokenPool, AddPoolReceiveContextNonContractValue, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  ApplyChainUpdates: 
    damlTypes.Choice<LockReleaseTokenPool, ApplyChainUpdates, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  ApplyTokenTransferFeeConfigUpdates: 
    damlTypes.Choice<LockReleaseTokenPool, ApplyTokenTransferFeeConfigUpdates, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  Archive: 
    damlTypes.Choice<LockReleaseTokenPool, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  CalculateFee: 
    damlTypes.Choice<LockReleaseTokenPool, CalculateFee, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  ClearPoolReceiveContext: 
    damlTypes.Choice<LockReleaseTokenPool, ClearPoolReceiveContext, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  GetFee: 
    damlTypes.Choice<LockReleaseTokenPool, GetFee, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.TokenPool.TokenPoolFeeQuote, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  GetRequiredCCVs: 
    damlTypes.Choice<LockReleaseTokenPool, GetRequiredCCVs, pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[], undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  LockOrBurn: 
    damlTypes.Choice<LockReleaseTokenPool, LockOrBurn, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.TokenPool.LockOrBurnResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  ReleaseFromTicket: 
    damlTypes.Choice<LockReleaseTokenPool, ReleaseFromTicket, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.TokenPool.ReleaseOrMintResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  RemovePoolReceiveContextValue: 
    damlTypes.Choice<LockReleaseTokenPool, RemovePoolReceiveContextValue, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  SetDynamicConfig: 
    damlTypes.Choice<LockReleaseTokenPool, SetDynamicConfig, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  SetRateLimitConfig: 
    damlTypes.Choice<LockReleaseTokenPool, SetRateLimitConfig, damlTypes.ContractId<pkg1377c6dda2856b806eba169ee7a1b545da7ad8a9416b53f90fdb389853c9f8e4.CCIP.RateLimiterV2.RateLimiter>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  SetRateLimiterReferences: 
    damlTypes.Choice<LockReleaseTokenPool, SetRateLimiterReferences, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  SetTransferTimeout: 
    damlTypes.Choice<LockReleaseTokenPool, SetTransferTimeout, damlTypes.ContractId<LockReleaseTokenPool>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  VerifyInboundMessage: 
    damlTypes.Choice<LockReleaseTokenPool, VerifyInboundMessage, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
  VerifyOutboundCCVs: 
    damlTypes.Choice<LockReleaseTokenPool, VerifyOutboundCCVs, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockReleaseTokenPool, undefined>>;
}
export declare const LockReleaseTokenPool:
  damlTypes.Template<LockReleaseTokenPool, undefined, '#ccip-lock-release-token-pool-v2:CCIP.LockReleaseTokenPoolV2:LockReleaseTokenPool'> &
  damlTypes.ToInterface<LockReleaseTokenPool, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.TokenPool.ITokenPool | pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LockReleaseTokenPoolInterface

export declare type LockReleaseTokenPoolDeps = {
  tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
}

export declare const LockReleaseTokenPoolDeps:
  damlTypes.Serializable<LockReleaseTokenPoolDeps>

export declare type ReleaseFromTicket = {
  tokenAdminRegistryCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenAdminRegistry>,
  tokenConfigCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenConfig>,
  rmnRemoteCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote>,
  tokenReceiveTicketCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.ITokenReceiveTicket>,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const ReleaseFromTicket:
  damlTypes.Serializable<ReleaseFromTicket>

export declare type RemovePoolReceiveContextValue = {
  contextKey: string,
}

export declare const RemovePoolReceiveContextValue:
  damlTypes.Serializable<RemovePoolReceiveContextValue>

export declare type SetDynamicConfig = {
  rateLimitAdmin: damlTypes.Optional<damlTypes.Party>,
}

export declare const SetDynamicConfig:
  damlTypes.Serializable<SetDynamicConfig>

export declare type SetRateLimitConfig = {
  caller: damlTypes.Party,
  rateLimiterCid: damlTypes.ContractId<pkg1377c6dda2856b806eba169ee7a1b545da7ad8a9416b53f90fdb389853c9f8e4.CCIP.RateLimiterV2.RateLimiter>,
  newIsEnabled: boolean,
  newCapacity: damlTypes.Numeric,
  newRate: damlTypes.Numeric,
}

export declare const SetRateLimitConfig:
  damlTypes.Serializable<SetRateLimitConfig>

export declare type SetRateLimiterReferences = {
  rateLimitConfigArgs: CCIP_LockReleaseTokenPoolV2Types.RateLimitConfigArgs[],
}

export declare const SetRateLimiterReferences:
  damlTypes.Serializable<SetRateLimiterReferences>

export declare type SetTransferTimeout = {
  newTransferTimeout: CCIP_LockReleaseTokenPoolV2Types.TransferTimeout,
}

export declare const SetTransferTimeout:
  damlTypes.Serializable<SetTransferTimeout>

export declare type VerifyInboundMessage = {
  tokenAdminRegistryCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenAdminRegistry>,
  tokenConfigCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenConfig>,
  executingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage>,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const VerifyInboundMessage:
  damlTypes.Serializable<VerifyInboundMessage>

export declare type VerifyOutboundCCVs = {
  tokenAdminRegistryCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenAdminRegistry>,
  tokenConfigCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.TokenAdminRegistry.ITokenConfig>,
  sendingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>,
  amount: damlTypes.Numeric,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const VerifyOutboundCCVs:
  damlTypes.Serializable<VerifyOutboundCCVs>
