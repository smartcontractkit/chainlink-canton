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

var pkg0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e = require('@chainlink/chainlink-canton-contracts/ccip-committee-verifier-v2-2.0.0');
var pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238 = require('@chainlink/chainlink-canton-contracts/ccip-runtime-v2-2.0.0');
var pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f = require('@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0');
var pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40 = require('@chainlink/chainlink-canton-contracts/ccip-rate-limiter-v2-2.0.0');
var pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 = require('@chainlink/chainlink-canton-contracts/mcms-api-1.0.0');
var pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb = require('@chainlink/chainlink-canton-contracts/ccip-codec-v2-2.0.0');
var pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b = require('@chainlink/chainlink-canton-contracts/splice-api-token-holding-v1-1.0.0');
var pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924 = require('@chainlink/chainlink-canton-contracts/ccip-core-v2-2.0.0');
var pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 = require('@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0');
var pkgb44e4259b6a243cdff7303f4ce28ce631d2e852dcf2837eba240d3c584cd9399 = require('@chainlink/chainlink-canton-contracts/ccip-burn-mint-token-pool-v2-2.0.0');
var pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8 = require('@chainlink/chainlink-canton-contracts/chainlink-api-2.0.0');
var pkgf0a7abce686b27e23ce7017b3b2159f7d7ecda430d101161fc5f10c67d5ac6d6 = require('@chainlink/chainlink-canton-contracts/ccip-lock-release-token-pool-v2-2.0.0');
var pkgf0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688 = require('@chainlink/chainlink-canton-contracts/ccip-executor-v2-2.0.0');
var pkgf83b5db150d68564a63d9a5760a8a49323424f58ebbea0e089fe79643a8f29c2 = require('@chainlink/chainlink-canton-contracts/link-2.0.0');

exports.CCIPFactory = damlTypes.assembleTemplate(
  {
    templateId: '#ccip-factory-v2:CCIP.FactoryV2:CCIPFactory',
    templateIdWithPackageId: '5b6cb3b6962cd5f91f2126c678b02e6d1dda54a5160d09e44edc087e36ed4930:CCIP.FactoryV2:CCIPFactory',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        instanceId: damlTypes.Text.decoder,
        owner: damlTypes.Party.decoder,
        mcmsParty: damlTypes.Party.decoder,
        usedInstanceIds: damlTypes.Map(damlTypes.Text, damlTypes.Bool).decoder,
        deployedContracts: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(damlTypes.Unit)).decoder,
        perPartyRouterFactoryDeployed: damlTypes.Bool.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        owner: damlTypes.Party.encode(__typed__.owner),
        mcmsParty: damlTypes.Party.encode(__typed__.mcmsParty),
        usedInstanceIds: damlTypes.Map(damlTypes.Text, damlTypes.Bool).encode(__typed__.usedInstanceIds),
        deployedContracts: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(damlTypes.Unit)).encode(__typed__.deployedContracts),
        perPartyRouterFactoryDeployed: damlTypes.Bool.encode(__typed__.perPartyRouterFactoryDeployed),
      };
    },
    Archive: {
      template: function () { return exports.CCIPFactory; },
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
    DeployBurnMintTokenPool: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployBurnMintTokenPool',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployBurnMintTokenPool.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployBurnMintTokenPool.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployCommitteeVerifier: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployCommitteeVerifier',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployCommitteeVerifier.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployCommitteeVerifier.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployExecutor: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployExecutor',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployExecutor.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployExecutor.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployFeeQuoter: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployFeeQuoter',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployFeeQuoter.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployFeeQuoter.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployGlobalConfig: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployGlobalConfig',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployGlobalConfig.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployGlobalConfig.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployLinkToken: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployLinkToken',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployLinkToken.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployLinkToken.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployLockReleaseTokenPool: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployLockReleaseTokenPool',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployLockReleaseTokenPool.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployLockReleaseTokenPool.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployOffRamp: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployOffRamp',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployOffRamp.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployOffRamp.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployOnRamp: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployOnRamp',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployOnRamp.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployOnRamp.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployPerPartyRouterFactory: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployPerPartyRouterFactory',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployPerPartyRouterFactory.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployPerPartyRouterFactory.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployRMNRemote: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployRMNRemote',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployRMNRemote.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployRMNRemote.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployRateLimiter: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployRateLimiter',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployRateLimiter.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployRateLimiter.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    DeployTokenAdminRegistry: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'DeployTokenAdminRegistry',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployTokenAdminRegistry.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DeployTokenAdminRegistry.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.DeployResult.decoder;
      }),
      resultEncode: function (__typed__) { return exports.DeployResult.encode(__typed__); },
    },
    GetFactoryState: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'GetFactoryState',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetFactoryState.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetFactoryState.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.FactoryState.decoder;
      }),
      resultEncode: function (__typed__) { return exports.FactoryState.encode(__typed__); },
    },
    SetOwnerToMCMS: {
      template: function () { return exports.CCIPFactory; },
      choiceName: 'SetOwnerToMCMS',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetOwnerToMCMS.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetOwnerToMCMS.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.CCIPFactory).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.CCIPFactory).encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.CCIPFactory, ['5b6cb3b6962cd5f91f2126c678b02e6d1dda54a5160d09e44edc087e36ed4930', '#ccip-factory-v2']);

exports.DeployBurnMintTokenPool = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkgb44e4259b6a243cdff7303f4ce28ce631d2e852dcf2837eba240d3c584cd9399.CCIP.BurnMintTokenPoolV2.BurnMintTokenPool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkgb44e4259b6a243cdff7303f4ce28ce631d2e852dcf2837eba240d3c584cd9399.CCIP.BurnMintTokenPoolV2.BurnMintTokenPool.encode(__typed__.contract),
    };
  },
};

exports.DeployBurnMintTokenPoolParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      poolOwner: damlTypes.Party.decoder,
      ccipOwner: damlTypes.Party.decoder,
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
      decimals: damlTypes.Int.decoder,
      rateLimitAdmin: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Party).decoder),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      transferTimeout: pkgb44e4259b6a243cdff7303f4ce28ce631d2e852dcf2837eba240d3c584cd9399.CCIP.BurnMintTokenPoolV2Types.TransferTimeout.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      poolOwner: damlTypes.Party.encode(__typed__.poolOwner),
      ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
      decimals: damlTypes.Int.encode(__typed__.decimals),
      rateLimitAdmin: damlTypes.Optional(damlTypes.Party).encode(__typed__.rateLimitAdmin),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.tokenAdminRegistry),
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.feeQuoter),
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
      poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.poolReceiveContext),
      transferTimeout: pkgb44e4259b6a243cdff7303f4ce28ce631d2e852dcf2837eba240d3c584cd9399.CCIP.BurnMintTokenPoolV2Types.TransferTimeout.encode(__typed__.transferTimeout),
    };
  },
};

exports.DeployCommitteeVerifier = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e.CCIP.CommitteeVerifierV2.CommitteeVerifier.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg0ad9d7694de68cd774fe399999722cf5fd8c514066125261e21b42431cfc962e.CCIP.CommitteeVerifierV2.CommitteeVerifier.encode(__typed__.contract),
    };
  },
};

exports.DeployCommitteeVerifierParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      owner: damlTypes.Party.decoder,
      ccipOwner: damlTypes.Party.decoder,
      versionTag: damlTypes.Text.decoder,
      allowListAdmin: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Party).decoder),
      messageSentObservers: damlTypes.List(damlTypes.Party).decoder,
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      storageLocations: damlTypes.List(damlTypes.Text).decoder,
      storageLocationsAdmin: damlTypes.Party.decoder,
      pendingStorageLocationsAdmin: damlTypes.Party.decoder,
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
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
      storageLocations: damlTypes.List(damlTypes.Text).encode(__typed__.storageLocations),
      storageLocationsAdmin: damlTypes.Party.encode(__typed__.storageLocationsAdmin),
      pendingStorageLocationsAdmin: damlTypes.Party.encode(__typed__.pendingStorageLocationsAdmin),
    };
  },
};

exports.DeployExecutor = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkgf0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688.CCIP.ExecutorV2.Executor.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkgf0fa05643b29b791f7928bcfb91f7bbf6c33622e91859fc3ef69c85482b4a688.CCIP.ExecutorV2.Executor.encode(__typed__.contract),
    };
  },
};

exports.DeployExecutorParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      owner: damlTypes.Party.decoder,
      maxCCVsPerMsg: damlTypes.Int.decoder,
      allowedFinalityConfig: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig.decoder,
      ccvAllowlistEnabled: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      owner: damlTypes.Party.encode(__typed__.owner),
      maxCCVsPerMsg: damlTypes.Int.encode(__typed__.maxCCVsPerMsg),
      allowedFinalityConfig: pkg6feabd6c3535eaaa23c820efbdaed64e15d733bdbfc292b88225888162774cfb.CCIP.CodecV2.FinalityConfig.FinalityConfig.encode(__typed__.allowedFinalityConfig),
      ccvAllowlistEnabled: damlTypes.Bool.encode(__typed__.ccvAllowlistEnabled),
    };
  },
};

exports.DeployFeeQuoter = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.FeeQuoter.FeeQuoter.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.FeeQuoter.FeeQuoter.encode(__typed__.contract),
    };
  },
};

exports.DeployFeeQuoterParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      linkTokenInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      linkTokenInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.linkTokenInstrumentId),
    };
  },
};

exports.DeployGlobalConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.GlobalConfig.GlobalConfig.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.GlobalConfig.GlobalConfig.encode(__typed__.contract),
    };
  },
};

exports.DeployGlobalConfigParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      chainSelector: damlTypes.Numeric(0).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      chainSelector: damlTypes.Numeric(0).encode(__typed__.chainSelector),
    };
  },
};

exports.DeployLinkToken = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkgf83b5db150d68564a63d9a5760a8a49323424f58ebbea0e089fe79643a8f29c2.Link.Token.LinkRegistry.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkgf83b5db150d68564a63d9a5760a8a49323424f58ebbea0e089fe79643a8f29c2.Link.Token.LinkRegistry.encode(__typed__.contract),
    };
  },
};

exports.DeployLinkTokenParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
    };
  },
};

exports.DeployLockReleaseTokenPool = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkgf0a7abce686b27e23ce7017b3b2159f7d7ecda430d101161fc5f10c67d5ac6d6.CCIP.LockReleaseTokenPoolV2.LockReleaseTokenPool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkgf0a7abce686b27e23ce7017b3b2159f7d7ecda430d101161fc5f10c67d5ac6d6.CCIP.LockReleaseTokenPoolV2.LockReleaseTokenPool.encode(__typed__.contract),
    };
  },
};

exports.DeployLockReleaseTokenPoolParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      poolOwner: damlTypes.Party.decoder,
      ccipOwner: damlTypes.Party.decoder,
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
      decimals: damlTypes.Int.decoder,
      rateLimitAdmin: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Party).decoder),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      transferTimeout: pkgf0a7abce686b27e23ce7017b3b2159f7d7ecda430d101161fc5f10c67d5ac6d6.CCIP.LockReleaseTokenPoolV2Types.TransferTimeout.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      poolOwner: damlTypes.Party.encode(__typed__.poolOwner),
      ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
      decimals: damlTypes.Int.encode(__typed__.decimals),
      rateLimitAdmin: damlTypes.Optional(damlTypes.Party).encode(__typed__.rateLimitAdmin),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.tokenAdminRegistry),
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.feeQuoter),
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
      poolReceiveContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.poolReceiveContext),
      transferTimeout: pkgf0a7abce686b27e23ce7017b3b2159f7d7ecda430d101161fc5f10c67d5ac6d6.CCIP.LockReleaseTokenPoolV2Types.TransferTimeout.encode(__typed__.transferTimeout),
    };
  },
};

exports.DeployOffRamp = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238.CCIP.RuntimeV2.OffRamp.OffRamp.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238.CCIP.RuntimeV2.OffRamp.OffRamp.encode(__typed__.contract),
    };
  },
};

exports.DeployOffRampParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.globalConfig),
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.tokenAdminRegistry),
    };
  },
};

exports.DeployOnRamp = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238.CCIP.RuntimeV2.OnRamp.OnRamp.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238.CCIP.RuntimeV2.OnRamp.OnRamp.encode(__typed__.contract),
    };
  },
};

exports.DeployOnRampParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      maxUSDCentsPerMsg: damlTypes.Numeric(0).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.globalConfig),
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.tokenAdminRegistry),
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.feeQuoter),
      maxUSDCentsPerMsg: damlTypes.Numeric(0).encode(__typed__.maxUSDCentsPerMsg),
    };
  },
};

exports.DeployPerPartyRouterFactory = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238.CCIP.RuntimeV2.PerPartyRouter.PerPartyRouterFactory.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg2235ebbb5f722d88a0e21aba9ced9abd0b637762768b8d38e17a4777fb8a3238.CCIP.RuntimeV2.PerPartyRouter.PerPartyRouterFactory.encode(__typed__.contract),
    };
  },
};

exports.DeployPerPartyRouterFactoryParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      onRamp: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      offRamp: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      onRamp: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.onRamp),
      offRamp: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.offRamp),
      globalConfig: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.globalConfig),
      tokenAdminRegistry: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.tokenAdminRegistry),
      feeQuoter: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.feeQuoter),
      rmnRemote: pkgb9f630bb75179b06f350030282a2276259016da59ba60677b3a8280d854dc2d8.Chainlink.InstanceAddress.RawInstanceAddress.encode(__typed__.rmnRemote),
    };
  },
};

exports.DeployRMNRemote = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.RMNRemote.RMNRemote.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.RMNRemote.RMNRemote.encode(__typed__.contract),
    };
  },
};

exports.DeployRMNRemoteParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      rmnOwner: damlTypes.Party.decoder,
      ccipOwner: damlTypes.Party.decoder,
      customObservers: damlTypes.List(damlTypes.Party).decoder,
      cursedSubjects: damlTypes.List(damlTypes.Text).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      rmnOwner: damlTypes.Party.encode(__typed__.rmnOwner),
      ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
      customObservers: damlTypes.List(damlTypes.Party).encode(__typed__.customObservers),
      cursedSubjects: damlTypes.List(damlTypes.Text).encode(__typed__.cursedSubjects),
    };
  },
};

exports.DeployRateLimiter = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimiter.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimiter.encode(__typed__.contract),
    };
  },
};

exports.DeployRateLimiterParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      poolInstanceId: damlTypes.Text.decoder,
      poolOwner: damlTypes.Party.decoder,
      remoteChainSelector: damlTypes.Numeric(0).decoder,
      direction: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimitDirection.decoder,
      mode: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimitMode.decoder,
      isEnabled: damlTypes.Bool.decoder,
      capacity: damlTypes.Numeric(0).decoder,
      rate: damlTypes.Numeric(0).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      poolInstanceId: damlTypes.Text.encode(__typed__.poolInstanceId),
      poolOwner: damlTypes.Party.encode(__typed__.poolOwner),
      remoteChainSelector: damlTypes.Numeric(0).encode(__typed__.remoteChainSelector),
      direction: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimitDirection.encode(__typed__.direction),
      mode: pkg62e9c12c62db5de2138722340d7ff597f9a07678aeba1e1fbf9ad121a9e89d40.CCIP.RateLimiterV2.RateLimitMode.encode(__typed__.mode),
      isEnabled: damlTypes.Bool.encode(__typed__.isEnabled),
      capacity: damlTypes.Numeric(0).encode(__typed__.capacity),
      rate: damlTypes.Numeric(0).encode(__typed__.rate),
    };
  },
};

exports.DeployResult = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      factoryCid: damlTypes.ContractId(exports.CCIPFactory).decoder,
      deployedCid: damlTypes.ContractId(damlTypes.Unit).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      factoryCid: damlTypes.ContractId(exports.CCIPFactory).encode(__typed__.factoryCid),
      deployedCid: damlTypes.ContractId(damlTypes.Unit).encode(__typed__.deployedCid),
    };
  },
};

exports.DeployTokenAdminRegistry = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.TokenAdminRegistry.TokenAdminRegistry.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      contract: pkg7f6d5a64b55cc88008c9664d166ffa8eee2998e17425106ea1b036baeed85924.CCIP.CoreV2.TokenAdminRegistry.TokenAdminRegistry.encode(__typed__.contract),
    };
  },
};

exports.DeployTokenAdminRegistryParams = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
    };
  },
};

exports.FactoryState = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      instanceId: damlTypes.Text.decoder,
      owner: damlTypes.Party.decoder,
      mcmsParty: damlTypes.Party.decoder,
      usedInstanceIds: damlTypes.Map(damlTypes.Text, damlTypes.Bool).decoder,
      deployedContracts: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(damlTypes.Unit)).decoder,
      perPartyRouterFactoryDeployed: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      instanceId: damlTypes.Text.encode(__typed__.instanceId),
      owner: damlTypes.Party.encode(__typed__.owner),
      mcmsParty: damlTypes.Party.encode(__typed__.mcmsParty),
      usedInstanceIds: damlTypes.Map(damlTypes.Text, damlTypes.Bool).encode(__typed__.usedInstanceIds),
      deployedContracts: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(damlTypes.Unit)).encode(__typed__.deployedContracts),
      perPartyRouterFactoryDeployed: damlTypes.Bool.encode(__typed__.perPartyRouterFactoryDeployed),
    };
  },
};

exports.GetFactoryState = {
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

exports.SetOwnerToMCMS = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};
