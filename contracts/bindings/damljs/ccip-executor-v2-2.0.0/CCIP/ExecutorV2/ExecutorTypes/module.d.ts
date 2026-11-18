// Generated from ../../../CCIP/ExecutorV2/ExecutorTypes/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb from '@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0';
import * as pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 from '@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0';

export declare type ApplyAllowedCCVUpdatesParams = {
  ccvsToRemove: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[],
  ccvsToAdd: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress[],
  ccvAllowlistEnabled: boolean,
}

export declare const ApplyAllowedCCVUpdatesParams:
  damlTypes.Serializable<ApplyAllowedCCVUpdatesParams>

export declare type ApplyDestChainUpdatesParams = {
  destChainSelectorsToRemove: damlTypes.Numeric[],
  destChainSelectorsToAdd: RemoteChainConfigArgs[],
}

export declare const ApplyDestChainUpdatesParams:
  damlTypes.Serializable<ApplyDestChainUpdatesParams>

export declare type DynamicConfig = {
  feeAggregator: damlTypes.Optional<damlTypes.Party>,
  allowedFinalityConfig: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig,
  ccvAllowlistEnabled: boolean,
}

export declare const DynamicConfig:
  damlTypes.Serializable<DynamicConfig>

export declare type RemoteChainConfig = {
  feeUSDCents: damlTypes.Numeric,
  enabled: boolean,
}

export declare const RemoteChainConfig:
  damlTypes.Serializable<RemoteChainConfig>

export declare type RemoteChainConfigArgs = {
  destChainSelector: damlTypes.Numeric,
  config: RemoteChainConfig,
}

export declare const RemoteChainConfigArgs:
  damlTypes.Serializable<RemoteChainConfigArgs>

export declare type SetDynamicConfigParams = {
  newDynamicConfig: DynamicConfig,
}

export declare const SetDynamicConfigParams:
  damlTypes.Serializable<SetDynamicConfigParams>
