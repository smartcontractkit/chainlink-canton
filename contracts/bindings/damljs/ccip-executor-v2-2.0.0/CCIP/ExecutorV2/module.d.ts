// Generated from ../../CCIP/ExecutorV2/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4 from '@chainlink/chainlink-canton-contracts/ccip-extension-api-v2-2.0.0';
import * as pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f from '@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0';
import * as pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 from '@chainlink/chainlink-canton-contracts/mcms-api-1.0.0';
import * as pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb from '@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0';
import * as pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58 from '@chainlink/chainlink-canton-contracts/ccip-api-v2-2.0.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';
import * as pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 from '@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0';

import * as CCIP_ExecutorV2_ExecutorTypes from '../../CCIP/ExecutorV2/ExecutorTypes/module';

export declare type ApplyAllowedCCVUpdates = {
  ccvsToRemove: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[],
  ccvsToAdd: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[],
  ccvAllowlistEnabled: boolean,
}

export declare const ApplyAllowedCCVUpdates:
  damlTypes.Serializable<ApplyAllowedCCVUpdates>

export declare type ApplyDestChainUpdates = {
  destChainSelectorsToRemove: damlTypes.Numeric[],
  destChainSelectorsToAdd: CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfigArgs[],
}

export declare const ApplyDestChainUpdates:
  damlTypes.Serializable<ApplyDestChainUpdates>

export declare type CalculateFee = {
  expectedExecutor: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  sendingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>,
  executorArgs: string,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const CalculateFee:
  damlTypes.Serializable<CalculateFee>

export declare type Executor = {
  instanceId: string,
  owner: damlTypes.Party,
  maxCCVsPerMsg: damlTypes.Int,
  dynamicConfig: CCIP_ExecutorV2_ExecutorTypes.DynamicConfig,
  allowedCCVs: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[],
  remoteChainConfigs: damlTypes.Map<damlTypes.Numeric, CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfig>,
}

export declare interface ExecutorInterface {
  ApplyAllowedCCVUpdates: 
    damlTypes.Choice<Executor, ApplyAllowedCCVUpdates, damlTypes.ContractId<Executor>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  ApplyDestChainUpdates: 
    damlTypes.Choice<Executor, ApplyDestChainUpdates, damlTypes.ContractId<Executor>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  Archive: 
    damlTypes.Choice<Executor, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  CalculateFee: 
    damlTypes.Choice<Executor, CalculateFee, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  GetAllowedCCVs: 
    damlTypes.Choice<Executor, GetAllowedCCVs, pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[], undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  GetAllowedFinalityConfig: 
    damlTypes.Choice<Executor, GetAllowedFinalityConfig, pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  GetDestChains: 
    damlTypes.Choice<Executor, GetDestChains, CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfigArgs[], undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  GetDynamicConfig: 
    damlTypes.Choice<Executor, GetDynamicConfig, CCIP_ExecutorV2_ExecutorTypes.DynamicConfig, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  GetFee: 
    damlTypes.Choice<Executor, GetFee, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.Executor.ExecutorFeeQuote, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  GetMaxCCVsPerMessage: 
    damlTypes.Choice<Executor, GetMaxCCVsPerMessage, damlTypes.Int, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
  SetDynamicConfig: 
    damlTypes.Choice<Executor, SetDynamicConfig, damlTypes.ContractId<Executor>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<Executor, undefined>>;
}
export declare const Executor:
  damlTypes.Template<Executor, undefined, '#ccip-executor-v2:CCIP.ExecutorV2:Executor'> &
  damlTypes.ToInterface<Executor, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.Executor.IExecutor | pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  ExecutorInterface

export declare type GetAllowedCCVs = {
  caller: damlTypes.Party,
}

export declare const GetAllowedCCVs:
  damlTypes.Serializable<GetAllowedCCVs>

export declare type GetAllowedFinalityConfig = {
  caller: damlTypes.Party,
}

export declare const GetAllowedFinalityConfig:
  damlTypes.Serializable<GetAllowedFinalityConfig>

export declare type GetDestChains = {
  caller: damlTypes.Party,
}

export declare const GetDestChains:
  damlTypes.Serializable<GetDestChains>

export declare type GetDynamicConfig = {
  caller: damlTypes.Party,
}

export declare const GetDynamicConfig:
  damlTypes.Serializable<GetDynamicConfig>

export declare type GetFee = {
  expectedExecutor: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
  destChainSelector: damlTypes.Numeric,
  requiredCCVs: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[],
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const GetFee:
  damlTypes.Serializable<GetFee>

export declare type GetMaxCCVsPerMessage = {
  caller: damlTypes.Party,
}

export declare const GetMaxCCVsPerMessage:
  damlTypes.Serializable<GetMaxCCVsPerMessage>

export declare type SetDynamicConfig = {
  newDynamicConfig: CCIP_ExecutorV2_ExecutorTypes.DynamicConfig,
}

export declare const SetDynamicConfig:
  damlTypes.Serializable<SetDynamicConfig>
