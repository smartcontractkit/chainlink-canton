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
var pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58 = require('@chainlink/chainlink-canton-contracts/ccip-api-v2-2.0.0');
var pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 = require('@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0');

var CCIP_CommitteeVerifierV2_CommitteeVerifierTypes = require('../../CCIP/CommitteeVerifierV2/CommitteeVerifierTypes/module');

exports.AcceptStorageLocationsAdmin = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.ApplyAllowListUpdates = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      allowListConfigArgsItems: damlTypes.List(CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.AllowListConfigArgs).decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      allowListConfigArgsItems: damlTypes.List(CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.AllowListConfigArgs).encode(__typed__.allowListConfigArgsItems),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.ApplyRemoteChainConfigUpdates = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      remoteChainConfigArgs: damlTypes.List(CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.RemoteChainConfigArgs).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      remoteChainConfigArgs: damlTypes.List(CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.RemoteChainConfigArgs).encode(__typed__.remoteChainConfigArgs),
    };
  },
};

exports.ApplySignatureConfigs = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      sourceChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).decoder,
      signatureConfigs: damlTypes.List(CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SignatureConfig).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      sourceChainSelectorsToRemove: damlTypes.List(damlTypes.Numeric(0)).encode(__typed__.sourceChainSelectorsToRemove),
      signatureConfigs: damlTypes.List(CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SignatureConfig).encode(__typed__.signatureConfigs),
    };
  },
};

exports.CalculateFee = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      sendingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).decoder,
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      sendingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).encode(__typed__.sendingMessageCid),
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.context),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.CommitteeVerifier = damlTypes.assembleTemplate(
  {
    templateId: '#ccip-committee-verifier-v2:CCIP.CommitteeVerifierV2:CommitteeVerifier',
    templateIdWithPackageId: '0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e:CCIP.CommitteeVerifierV2:CommitteeVerifier',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        instanceId: damlTypes.Text.decoder,
        owner: damlTypes.Party.decoder,
        ccipOwner: damlTypes.Party.decoder,
        versionTag: damlTypes.Text.decoder,
        allowListAdmin: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Party).decoder),
        messageSentObservers: damlTypes.List(damlTypes.Party).decoder,
        storageLocations: damlTypes.List(damlTypes.Text).decoder,
        storageLocationsAdmin: damlTypes.Party.decoder,
        pendingStorageLocationsAdmin: damlTypes.Party.decoder,
        remoteChainConfigs: damlTypes.Map(damlTypes.Numeric(0), CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.RemoteChainConfig).decoder,
        signerConfigs: damlTypes.Map(damlTypes.Numeric(0), CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SignatureConfig).decoder,
        deps: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.CommitteeVerifierDeps.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        owner: damlTypes.Party.encode(__typed__.owner),
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        versionTag: damlTypes.Text.encode(__typed__.versionTag),
        allowListAdmin: damlTypes.Optional(damlTypes.Party).encode(__typed__.allowListAdmin),
        messageSentObservers: damlTypes.List(damlTypes.Party).encode(__typed__.messageSentObservers),
        storageLocations: damlTypes.List(damlTypes.Text).encode(__typed__.storageLocations),
        storageLocationsAdmin: damlTypes.Party.encode(__typed__.storageLocationsAdmin),
        pendingStorageLocationsAdmin: damlTypes.Party.encode(__typed__.pendingStorageLocationsAdmin),
        remoteChainConfigs: damlTypes.Map(damlTypes.Numeric(0), CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.RemoteChainConfig).encode(__typed__.remoteChainConfigs),
        signerConfigs: damlTypes.Map(damlTypes.Numeric(0), CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SignatureConfig).encode(__typed__.signerConfigs),
        deps: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.CommitteeVerifierDeps.encode(__typed__.deps),
      };
    },
    AcceptStorageLocationsAdmin: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'AcceptStorageLocationsAdmin',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.AcceptStorageLocationsAdmin.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.AcceptStorageLocationsAdmin.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    ApplyAllowListUpdates: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'ApplyAllowListUpdates',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApplyAllowListUpdates.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApplyAllowListUpdates.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    ApplyRemoteChainConfigUpdates: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'ApplyRemoteChainConfigUpdates',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApplyRemoteChainConfigUpdates.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApplyRemoteChainConfigUpdates.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    ApplySignatureConfigs: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'ApplySignatureConfigs',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApplySignatureConfigs.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApplySignatureConfigs.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    Archive: {
      template: function () { return exports.CommitteeVerifier; },
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
      template: function () { return exports.CommitteeVerifier; },
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
    ForwardToVerifier: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'ForwardToVerifier',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ForwardToVerifier.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ForwardToVerifier.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).encode(__typed__); },
    },
    GetFee: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'GetFee',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetFee.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetFee.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.CrossChainVerifier.CrossChainVerifierFeeQuote.decoder;
      }),
      resultEncode: function (__typed__) { return pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.CrossChainVerifier.CrossChainVerifierFeeQuote.encode(__typed__); },
    },
    SetDeps: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'SetDeps',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetDeps.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetDeps.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    SetDynamicConfig: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'SetDynamicConfig',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetDynamicConfig.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetDynamicConfig.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    TransferStorageLocationsAdmin: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'TransferStorageLocationsAdmin',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.TransferStorageLocationsAdmin.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.TransferStorageLocationsAdmin.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    UpdateStorageLocations: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'UpdateStorageLocations',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.UpdateStorageLocations.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.UpdateStorageLocations.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CommitteeVerifier).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CommitteeVerifier).encode(__typed__); },
    },
    VerifyMessage: {
      template: function () { return exports.CommitteeVerifier; },
      choiceName: 'VerifyMessage',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.VerifyMessage.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.VerifyMessage.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage).encode(__typed__); },
    },
  },
  pkg289011bdbefe42c7dbea0a4f101127095e8b5f5281d45c84f6eca06de11689a4.CCIP.InterfacesV2.CrossChainVerifier.ICrossChainVerifier,
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.CommitteeVerifier, ['0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e', '#ccip-committee-verifier-v2']);

exports.ForwardToVerifier = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      rmnRemoteCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote).decoder,
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      sendingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).decoder,
      verifierArgs: damlTypes.Text.decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      rmnRemoteCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote).encode(__typed__.rmnRemoteCid),
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.context),
      sendingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.SendingMessage.ISendingMessage).encode(__typed__.sendingMessageCid),
      verifierArgs: damlTypes.Text.encode(__typed__.verifierArgs),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.GetFee = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      destChainSelector: damlTypes.Numeric(0).decoder,
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      destChainSelector: damlTypes.Numeric(0).encode(__typed__.destChainSelector),
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.context),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};

exports.SetDeps = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      newDeps: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SetDepsParams.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      newDeps: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.SetDepsParams.encode(__typed__.newDeps),
    };
  },
};

exports.SetDynamicConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      dynamicConfig: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.DynamicConfig.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      dynamicConfig: CCIP_CommitteeVerifierV2_CommitteeVerifierTypes.DynamicConfig.encode(__typed__.dynamicConfig),
    };
  },
};

exports.TransferStorageLocationsAdmin = {
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

exports.UpdateStorageLocations = {
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

exports.VerifyMessage = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      rmnRemoteCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote).decoder,
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      executingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage).decoder,
      verifierResults: damlTypes.Text.decoder,
      caller: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      rmnRemoteCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.RMNRemote.IRMNRemote).encode(__typed__.rmnRemoteCid),
      context: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.context),
      executingMessageCid: damlTypes.ContractId(pkg7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58.CCIP.APIV2.ExecutingMessage.IExecutingMessage).encode(__typed__.executingMessageCid),
      verifierResults: damlTypes.Text.encode(__typed__.verifierResults),
      caller: damlTypes.Party.encode(__typed__.caller),
    };
  },
};
