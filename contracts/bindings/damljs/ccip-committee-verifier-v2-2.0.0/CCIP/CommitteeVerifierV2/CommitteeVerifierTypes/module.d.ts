// Generated from ../../../CCIP/CommitteeVerifierV2/CommitteeVerifierTypes/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 from '@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0';

export declare type AllowListConfigArgs = {
  destChainSelector: damlTypes.Numeric,
  allowListEnabled: boolean,
  addedAllowListedSenders: damlTypes.Party[],
  removedAllowListedSenders: damlTypes.Party[],
}

export declare const AllowListConfigArgs:
  damlTypes.Serializable<AllowListConfigArgs>

export declare type ApplyAllowListUpdatesParams = {
  allowListConfigArgs: AllowListConfigArgs[],
  caller: damlTypes.Party,
}

export declare const ApplyAllowListUpdatesParams:
  damlTypes.Serializable<ApplyAllowListUpdatesParams>

export declare type ApplyRemoteChainConfigUpdatesParams = {
  remoteChainConfigArgs: RemoteChainConfigArgs[],
}

export declare const ApplyRemoteChainConfigUpdatesParams:
  damlTypes.Serializable<ApplyRemoteChainConfigUpdatesParams>

export declare type ApplySignatureConfigsParams = {
  sourceChainSelectorsToRemove: damlTypes.Numeric[],
  signatureConfigs: SignatureConfig[],
}

export declare const ApplySignatureConfigsParams:
  damlTypes.Serializable<ApplySignatureConfigsParams>

export declare type CommitteeVerifierDeps = {
  rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress,
}

export declare const CommitteeVerifierDeps:
  damlTypes.Serializable<CommitteeVerifierDeps>

export declare type DynamicConfig = {
  allowListAdmin: damlTypes.Optional<damlTypes.Party>,
  messageSentObservers: damlTypes.Party[],
}

export declare const DynamicConfig:
  damlTypes.Serializable<DynamicConfig>

export declare type RemoteChainConfig = {
  feeUSDCents: damlTypes.Numeric,
  gasForVerification: damlTypes.Int,
  payloadSizeBytes: damlTypes.Int,
  allowListEnabled: boolean,
  allowedSendersList: damlTypes.Party[],
}

export declare const RemoteChainConfig:
  damlTypes.Serializable<RemoteChainConfig>

export declare type RemoteChainConfigArgs = {
  remoteChainSelector: damlTypes.Numeric,
  feeUSDCents: damlTypes.Numeric,
  gasForVerification: damlTypes.Int,
  payloadSizeBytes: damlTypes.Int,
  allowListEnabled: boolean,
}

export declare const RemoteChainConfigArgs:
  damlTypes.Serializable<RemoteChainConfigArgs>

export declare type SetDepsParams = {
  rmnRemote: damlTypes.Optional<pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress>,
}

export declare const SetDepsParams:
  damlTypes.Serializable<SetDepsParams>

export declare type SetDynamicConfigParams = {
  dynamicConfig: DynamicConfig,
}

export declare const SetDynamicConfigParams:
  damlTypes.Serializable<SetDynamicConfigParams>

export declare type SignatureConfig = {
  sourceChainSelector: damlTypes.Numeric,
  threshold: damlTypes.Int,
  signerKeys: string[],
}

export declare const SignatureConfig:
  damlTypes.Serializable<SignatureConfig>

export declare type TransferStorageLocationsAdminParams = {
  newAdmin: damlTypes.Party,
}

export declare const TransferStorageLocationsAdminParams:
  damlTypes.Serializable<TransferStorageLocationsAdminParams>

export declare type UpdateStorageLocationsParams = {
  newLocations: string[],
}

export declare const UpdateStorageLocationsParams:
  damlTypes.Serializable<UpdateStorageLocationsParams>
