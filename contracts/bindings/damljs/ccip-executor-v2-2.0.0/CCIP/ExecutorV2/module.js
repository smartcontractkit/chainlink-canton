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

var pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4 = require('@chainlink/chainlink-canton-contracts/ccip-extension-api-v2-2.0.0');
var pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f = require('@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0');
var pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 = require('@chainlink/chainlink-canton-contracts/mcms-api-1.0.0');
var pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb = require('@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0');
var pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58 = require('@chainlink/chainlink-canton-contracts/ccip-api-v2-2.0.0');
var pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 = require('@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0');
var pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 = require('@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0');

var CCIP_ExecutorV2_ExecutorTypes = require('../../CCIP/ExecutorV2/ExecutorTypes/module');

exports.ApplyAllowedCCVUpdates = {
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

exports.ApplyDestChainUpdates = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      destChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).decoder,
      destChainSelectorsToAdd: damlTypes.List(CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfigArgs).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      destChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).encode(__typed__.destChainSelectorsToRemove),
      destChainSelectorsToAdd: damlTypes.List(CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfigArgs).encode(__typed__.destChainSelectorsToAdd),
    };
  },
};

exports.CalculateFee = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      expectedExecutor: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      sendingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).decoder,
      executorArgs: damlTypes.Text.decoder,
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      expectedExecutor: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.expectedExecutor),
      sendingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).encode(__typed__.sendingMessageCid),
      executorArgs: damlTypes.Text.encode(__typed__.executorArgs),
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.context),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.Executor = damlTypes.assembleTemplate(
  {
    templateId: '#ccip-executor-v2:CCIP.ExecutorV2:Executor',
    templateIdWithPackageId: 'f0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688:CCIP.ExecutorV2:Executor',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        instanceId: damlTypes.Text.decoder,
        owner: damlTypes.Party.decoder,
        maxCCVsPerMsg: damlTypes.Int.decoder,
        dynamicConfig: CCIP_ExecutorV2_ExecutorTypes.DynamicConfig.decoder,
        allowedCCVs: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).decoder,
        remoteChainConfigs: damlTypes.Map(damlTypes.Numeric(0), CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfig).decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        owner: damlTypes.Party.encode(__typed__.owner),
        maxCCVsPerMsg: damlTypes.Int.encode(__typed__.maxCCVsPerMsg),
        dynamicConfig: CCIP_ExecutorV2_ExecutorTypes.DynamicConfig.encode(__typed__.dynamicConfig),
        allowedCCVs: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).encode(__typed__.allowedCCVs),
        remoteChainConfigs: damlTypes.Map(damlTypes.Numeric(0), CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfig).encode(__typed__.remoteChainConfigs),
      };
    },
    ApplyAllowedCCVUpdates: {
      template: function () { return exports.Executor; },
      choiceName: 'ApplyAllowedCCVUpdates',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApplyAllowedCCVUpdates.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApplyAllowedCCVUpdates.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.Executor).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.Executor).encode(__typed__); },
    },
    ApplyDestChainUpdates: {
      template: function () { return exports.Executor; },
      choiceName: 'ApplyDestChainUpdates',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApplyDestChainUpdates.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApplyDestChainUpdates.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.Executor).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.Executor).encode(__typed__); },
    },
    Archive: {
      template: function () { return exports.Executor; },
      choiceName: 'Archive',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive.decoder;
      }),
      argumentEncode: function (__typed__) { return pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    CalculateFee: {
      template: function () { return exports.Executor; },
      choiceName: 'CalculateFee',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CalculateFee.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CalculateFee.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).encode(__typed__); },
    },
    GetAllowedCCVs: {
      template: function () { return exports.Executor; },
      choiceName: 'GetAllowedCCVs',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetAllowedCCVs.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetAllowedCCVs.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).encode(__typed__); },
    },
    GetAllowedFinalityConfig: {
      template: function () { return exports.Executor; },
      choiceName: 'GetAllowedFinalityConfig',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetAllowedFinalityConfig.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetAllowedFinalityConfig.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig.decoder;
      }),
      resultEncode: function (__typed__) { return pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig.encode(__typed__); },
    },
    GetDestChains: {
      template: function () { return exports.Executor; },
      choiceName: 'GetDestChains',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetDestChains.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetDestChains.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.List(CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfigArgs).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.List(CCIP_ExecutorV2_ExecutorTypes.RemoteChainConfigArgs).encode(__typed__); },
    },
    GetDynamicConfig: {
      template: function () { return exports.Executor; },
      choiceName: 'GetDynamicConfig',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetDynamicConfig.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetDynamicConfig.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return CCIP_ExecutorV2_ExecutorTypes.DynamicConfig.decoder;
      }),
      resultEncode: function (__typed__) { return CCIP_ExecutorV2_ExecutorTypes.DynamicConfig.encode(__typed__); },
    },
    GetFee: {
      template: function () { return exports.Executor; },
      choiceName: 'GetFee',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetFee.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetFee.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.Executor.ExecutorFeeQuote.decoder;
      }),
      resultEncode: function (__typed__) { return pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.Executor.ExecutorFeeQuote.encode(__typed__); },
    },
    GetMaxCCVsPerMessage: {
      template: function () { return exports.Executor; },
      choiceName: 'GetMaxCCVsPerMessage',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetMaxCCVsPerMessage.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetMaxCCVsPerMessage.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Int.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Int.encode(__typed__); },
    },
    SetDynamicConfig: {
      template: function () { return exports.Executor; },
      choiceName: 'SetDynamicConfig',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetDynamicConfig.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetDynamicConfig.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.Executor).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.Executor).encode(__typed__); },
    },
  },
  pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.Executor.IExecutor,
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.Executor, ['f0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688', '#ccip-executor-v2']);

exports.GetAllowedCCVs = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.GetAllowedFinalityConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.GetDestChains = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.GetDynamicConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.GetFee = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      expectedExecutor: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      destChainSelector: damlTypes.Numeric(0).decoder,
      requiredCCVs: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).decoder,
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      expectedExecutor: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.expectedExecutor),
      destChainSelector: damlTypes.Numeric(0).encode(__typed__.destChainSelector),
      requiredCCVs: damlTypes.List(pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress).encode(__typed__.requiredCCVs),
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.context),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.GetMaxCCVsPerMessage = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.SetDynamicConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      newDynamicConfig: CCIP_ExecutorV2_ExecutorTypes.DynamicConfig.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      newDynamicConfig: CCIP_ExecutorV2_ExecutorTypes.DynamicConfig.encode(__typed__.newDynamicConfig),
    };
  },
};
