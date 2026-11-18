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

var pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 = require('@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0');

exports.AllowListConfigArgs = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      destChainSelector: damlTypes.Numeric(0).decoder,
      allowListEnabled: damlTypes.Bool.decoder,
      addedAllowListedSenders: damlTypes.List(damlTypes.Party).decoder,
      removedAllowListedSenders: damlTypes.List(damlTypes.Party).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      destChainSelector: damlTypes.Numeric(0).encode(__typed__.destChainSelector),
      allowListEnabled: damlTypes.Bool.encode(__typed__.allowListEnabled),
      addedAllowListedSenders: damlTypes.List(damlTypes.Party).encode(__typed__.addedAllowListedSenders),
      removedAllowListedSenders: damlTypes.List(damlTypes.Party).encode(__typed__.removedAllowListedSenders),
    };
  },
};

exports.ApplyAllowListUpdatesParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      allowListConfigArgs: damlTypes.List(exports.AllowListConfigArgs).decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      allowListConfigArgs: damlTypes.List(exports.AllowListConfigArgs).encode(__typed__.allowListConfigArgs),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.ApplyRemoteChainConfigUpdatesParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      remoteChainConfigArgs: damlTypes.List(exports.RemoteChainConfigArgs).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      remoteChainConfigArgs: damlTypes.List(exports.RemoteChainConfigArgs).encode(__typed__.remoteChainConfigArgs),
    };
  },
};

exports.ApplySignatureConfigsParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      sourceChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).decoder,
      signatureConfigs: damlTypes.List(exports.SignatureConfig).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      sourceChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).encode(__typed__.sourceChainSelectorsToRemove),
      signatureConfigs: damlTypes.List(exports.SignatureConfig).encode(__typed__.signatureConfigs),
    };
  },
};

exports.CommitteeVerifierDeps = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
    };
  },
};

exports.DynamicConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      allowListAdmin: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Party).decoder),
      messageSentObservers: damlTypes.List(damlTypes.Party).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      allowListAdmin: damlTypes.Optional(damlTypes.Party).encode(__typed__.allowListAdmin),
      messageSentObservers: damlTypes.List(damlTypes.Party).encode(__typed__.messageSentObservers),
    };
  },
};

exports.RemoteChainConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      feeUSDCents: damlTypes.Numeric(0).decoder,
      gasForVerification: damlTypes.Int.decoder,
      payloadSizeBytes: damlTypes.Int.decoder,
      allowListEnabled: damlTypes.Bool.decoder,
      allowedSendersList: damlTypes.List(damlTypes.Party).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      feeUSDCents: damlTypes.Numeric(0).encode(__typed__.feeUSDCents),
      gasForVerification: damlTypes.Int.encode(__typed__.gasForVerification),
      payloadSizeBytes: damlTypes.Int.encode(__typed__.payloadSizeBytes),
      allowListEnabled: damlTypes.Bool.encode(__typed__.allowListEnabled),
      allowedSendersList: damlTypes.List(damlTypes.Party).encode(__typed__.allowedSendersList),
    };
  },
};

exports.RemoteChainConfigArgs = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      remoteChainSelector: damlTypes.Numeric(0).decoder,
      feeUSDCents: damlTypes.Numeric(0).decoder,
      gasForVerification: damlTypes.Int.decoder,
      payloadSizeBytes: damlTypes.Int.decoder,
      allowListEnabled: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      remoteChainSelector: damlTypes.Numeric(0).encode(__typed__.remoteChainSelector),
      feeUSDCents: damlTypes.Numeric(0).encode(__typed__.feeUSDCents),
      gasForVerification: damlTypes.Int.encode(__typed__.gasForVerification),
      payloadSizeBytes: damlTypes.Int.encode(__typed__.payloadSizeBytes),
      allowListEnabled: damlTypes.Bool.encode(__typed__.allowListEnabled),
    };
  },
};

exports.SetDepsParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      rmnRemote: jtv.Decoder.withDefault(null, damlTypes.Optional(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).decoder),
    });
  }),
  encode: function (__typed__) {
    return {
      rmnRemote: damlTypes.Optional(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).encode(__typed__.rmnRemote),
    };
  },
};

exports.SetDynamicConfigParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      dynamicConfig: exports.DynamicConfig.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      dynamicConfig: exports.DynamicConfig.encode(__typed__.dynamicConfig),
    };
  },
};

exports.SignatureConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      sourceChainSelector: damlTypes.Numeric(0).decoder,
      threshold: damlTypes.Int.decoder,
      signerKeys: damlTypes.List(damlTypes.Text).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      sourceChainSelector: damlTypes.Numeric(0).encode(__typed__.sourceChainSelector),
      threshold: damlTypes.Int.encode(__typed__.threshold),
      signerKeys: damlTypes.List(damlTypes.Text).encode(__typed__.signerKeys),
    };
  },
};

exports.TransferStorageLocationsAdminParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      newAdmin: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      newAdmin: damlTypes.Party.encode(__typed__.newAdmin),
    };
  },
};

exports.UpdateStorageLocationsParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      newLocations: damlTypes.List(damlTypes.Text).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      newLocations: damlTypes.List(damlTypes.Text).encode(__typed__.newLocations),
    };
  },
};
