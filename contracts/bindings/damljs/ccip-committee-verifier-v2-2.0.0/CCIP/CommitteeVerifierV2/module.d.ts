// Generated from ../../CCIP/CommitteeVerifierV2/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4 from '@chainlink/chainlink-canton-contracts/ccip-extension-api-v2-2.0.0';
import * as pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f from '@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0';
import * as pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 from '@chainlink/chainlink-canton-contracts/mcms-api-1.0.0';
import * as pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58 from '@chainlink/chainlink-canton-contracts/ccip-api-v2-2.0.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';

import * as CCIP_CommitteeVerifierV2_CommitteeVerifierTypes from '../../CCIP/CommitteeVerifierV2/CommitteeVerifierTypes/module';

export declare type AcceptStorageLocationsAdmin = {
}

export declare const AcceptStorageLocationsAdmin:
  damlTypes.Serializable<AcceptStorageLocationsAdmin>

export declare type ApplyAllowListUpdates = {
  allowListConfigArgsItems: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.AllowListConfigArgs[],
  caller: damlTypes.Party,
}

export declare const ApplyAllowListUpdates:
  damlTypes.Serializable<ApplyAllowListUpdates>

export declare type ApplyRemoteChainConfigUpdates = {
  remoteChainConfigArgs: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.RemoteChainConfigArgs[],
}

export declare const ApplyRemoteChainConfigUpdates:
  damlTypes.Serializable<ApplyRemoteChainConfigUpdates>

export declare type ApplySignatureConfigs = {
  sourceChainSelectorsToRemove: damlTypes.Numeric[],
  signatureConfigs: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SignatureConfig[],
}

export declare const ApplySignatureConfigs:
  damlTypes.Serializable<ApplySignatureConfigs>

export declare type CalculateFee = {
  sendingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const CalculateFee:
  damlTypes.Serializable<CalculateFee>

export declare type CommitteeVerifier = {
  instanceId: string,
  owner: damlTypes.Party,
  ccipOwner: damlTypes.Party,
  versionTag: string,
  allowListAdmin: damlTypes.Optional<damlTypes.Party>,
  messageSentObservers: damlTypes.Party[],
  storageLocations: string[],
  storageLocationsAdmin: damlTypes.Party,
  pendingStorageLocationsAdmin: damlTypes.Party,
  remoteChainConfigs: damlTypes.Map<damlTypes.Numeric, CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.RemoteChainConfig>,
  signerConfigs: damlTypes.Map<damlTypes.Numeric, CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SignatureConfig>,
  deps: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.CommitteeVerifierDeps,
}

export declare interface CommitteeVerifierInterface {
  AcceptStorageLocationsAdmin: 
    damlTypes.Choice<CommitteeVerifier, AcceptStorageLocationsAdmin, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  ApplyAllowListUpdates: 
    damlTypes.Choice<CommitteeVerifier, ApplyAllowListUpdates, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  ApplyRemoteChainConfigUpdates: 
    damlTypes.Choice<CommitteeVerifier, ApplyRemoteChainConfigUpdates, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  ApplySignatureConfigs: 
    damlTypes.Choice<CommitteeVerifier, ApplySignatureConfigs, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  Archive: 
    damlTypes.Choice<CommitteeVerifier, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  CalculateFee: 
    damlTypes.Choice<CommitteeVerifier, CalculateFee, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  ForwardToVerifier: 
    damlTypes.Choice<CommitteeVerifier, ForwardToVerifier, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  GetFee: 
    damlTypes.Choice<CommitteeVerifier, GetFee, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.CrossChainVerifier.CrossChainVerifierFeeQuote, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  SetDeps: 
    damlTypes.Choice<CommitteeVerifier, SetDeps, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  SetDynamicConfig: 
    damlTypes.Choice<CommitteeVerifier, SetDynamicConfig, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  TransferStorageLocationsAdmin: 
    damlTypes.Choice<CommitteeVerifier, TransferStorageLocationsAdmin, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  UpdateStorageLocations: 
    damlTypes.Choice<CommitteeVerifier, UpdateStorageLocations, damlTypes.ContractId<CommitteeVerifier>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
  VerifyMessage: 
    damlTypes.Choice<CommitteeVerifier, VerifyMessage, damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<CommitteeVerifier, undefined>>;
}
export declare const CommitteeVerifier:
  damlTypes.Template<CommitteeVerifier, undefined, '#ccip-committee-verifier-v2:CCIP.CommitteeVerifierV2:CommitteeVerifier'> &
  damlTypes.ToInterface<CommitteeVerifier, pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.CrossChainVerifier.ICrossChainVerifier | pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  CommitteeVerifierInterface

export declare type ForwardToVerifier = {
  rmnRemoteCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote>,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  sendingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage>,
  verifierArgs: string,
  caller: damlTypes.Party,
}

export declare const ForwardToVerifier:
  damlTypes.Serializable<ForwardToVerifier>

export declare type GetFee = {
  destChainSelector: damlTypes.Numeric,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  caller: damlTypes.Party,
}

export declare const GetFee:
  damlTypes.Serializable<GetFee>

export declare type SetDeps = {
  newDeps: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SetDepsParams,
}

export declare const SetDeps:
  damlTypes.Serializable<SetDeps>

export declare type SetDynamicConfig = {
  dynamicConfig: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.DynamicConfig,
}

export declare const SetDynamicConfig:
  damlTypes.Serializable<SetDynamicConfig>

export declare type TransferStorageLocationsAdmin = {
  newAdmin: damlTypes.Party,
}

export declare const TransferStorageLocationsAdmin:
  damlTypes.Serializable<TransferStorageLocationsAdmin>

export declare type UpdateStorageLocations = {
  newLocations: string[],
}

export declare const UpdateStorageLocations:
  damlTypes.Serializable<UpdateStorageLocations>

export declare type VerifyMessage = {
  rmnRemoteCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote>,
  context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  executingMessageCid: damlTypes.ContractId<pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage>,
  verifierResults: string,
  caller: damlTypes.Party,
}

export declare const VerifyMessage:
  damlTypes.Serializable<VerifyMessage>
