"use strict";
/* eslint-disable-next-line no-unused-vars */
function __export(m) {
/* eslint-disable-next-line no-prototype-builtins */
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });

/* eslint-disable-next-line no-unused-vars */
var jtv = require('@mojotech/json-type-validation');
/* eslint-disable-next-line no-unused-vars */
var damlTypes = require('@daml/types');

var pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb = require('@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0');
var pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 = require('@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0');

exports.ApplyAllowedCCVUpdatesParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      ccvsToRemove: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).decoder,
      ccvsToAdd: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).decoder,
      ccvAllowlistEnabled: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      ccvsToRemove: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).encode(__typed__.ccvsToRemove),
      ccvsToAdd: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).encode(__typed__.ccvsToAdd),
      ccvAllowlistEnabled: damlTypes.Bool.encode(__typed__.ccvAllowlistEnabled),
    };
  },
};

exports.ApplyDestChainUpdatesParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      destChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).decoder,
      destChainSelectorsToAdd: damlTypes.List(exports.RemoteChainConfigArgs).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      destChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).encode(__typed__.destChainSelectorsToRemove),
      destChainSelectorsToAdd: damlTypes.List(exports.RemoteChainConfigArgs).encode(__typed__.destChainSelectorsToAdd),
    };
  },
};

exports.DynamicConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      feeAggregator: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Party).decoder),
      allowedFinalityConfig: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig.decoder,
      ccvAllowlistEnabled: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      feeAggregator: damlTypes.Optional(damlTypes.Party).encode(__typed__.feeAggregator),
      allowedFinalityConfig: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig.encode(__typed__.allowedFinalityConfig),
      ccvAllowlistEnabled: damlTypes.Bool.encode(__typed__.ccvAllowlistEnabled),
    };
  },
};

exports.RemoteChainConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      feeUSDCents: damlTypes.Numeric(0).decoder,
      enabled: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      feeUSDCents: damlTypes.Numeric(0).encode(__typed__.feeUSDCents),
      enabled: damlTypes.Bool.encode(__typed__.enabled),
    };
  },
};

exports.RemoteChainConfigArgs = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      destChainSelector: damlTypes.Numeric(0).decoder,
      config: exports.RemoteChainConfig.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      destChainSelector: damlTypes.Numeric(0).encode(__typed__.destChainSelector),
      config: exports.RemoteChainConfig.encode(__typed__.config),
    };
  },
};

exports.SetDynamicConfigParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      newDynamicConfig: exports.DynamicConfig.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      newDynamicConfig: exports.DynamicConfig.encode(__typed__.newDynamicConfig),
    };
  },
};
