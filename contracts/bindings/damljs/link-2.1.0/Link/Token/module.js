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

var pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f = require('@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0');
var pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281 = require('@chainlink/chainlink-canton-contracts/splice-api-token-transfer-instruction-v1-1.0.0');
var pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b = require('@chainlink/chainlink-canton-contracts/splice-api-token-holding-v1-1.0.0');
var pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e = require('@chainlink/chainlink-canton-contracts/splice-api-token-burn-mint-v1-1.0.0');
var pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 = require('@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0');

exports.Cancel = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      actor: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      actor: damlTypes.Party.encode(__typed__.actor),
    };
  },
};

exports.CreateTransferPreapproval = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      receiver: damlTypes.Party.decoder,
      existingPreapprovalCid: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.ContractId(exports.LinkTransferPreapproval)).decoder),
    });
  }),
  encode: function (__typed__) {
    return {
      receiver: damlTypes.Party.encode(__typed__.receiver),
      existingPreapprovalCid: damlTypes.Optional(damlTypes.ContractId(exports.LinkTransferPreapproval)).encode(__typed__.existingPreapprovalCid),
    };
  },
};

exports.LinkHolding = damlTypes.assembleTemplate(
  {
    templateId: '#link:Link.Token:LinkHolding',
    templateIdWithPackageId: 'fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60:Link.Token:LinkHolding',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        holdingOwner: damlTypes.Party.decoder,
        holdingAdmin: damlTypes.Party.decoder,
        holdingInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        holdingAmount: damlTypes.Numeric(10).decoder,
        meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        holdingOwner: damlTypes.Party.encode(__typed__.holdingOwner),
        holdingAdmin: damlTypes.Party.encode(__typed__.holdingAdmin),
        holdingInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.holdingInstrumentId),
        holdingAmount: damlTypes.Numeric(10).encode(__typed__.holdingAmount),
        meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
      };
    },
    Archive: {
      template: function () { return exports.LinkHolding; },
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
  },
  pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding,
);

damlTypes.registerTemplate(exports.LinkHolding, ['fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60', '#link']);

exports.LinkRegistry = damlTypes.assembleTemplate(
  {
    templateId: '#link:Link.Token:LinkRegistry',
    templateIdWithPackageId: 'fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60:Link.Token:LinkRegistry',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        registryAdmin: damlTypes.Party.decoder,
        registryInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        instanceId: damlTypes.Text.decoder,
        registryMeta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        registryAdmin: damlTypes.Party.encode(__typed__.registryAdmin),
        registryInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.registryInstrumentId),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        registryMeta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.registryMeta),
      };
    },
    Archive: {
      template: function () { return exports.LinkRegistry; },
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
    CreateTransferPreapproval: {
      template: function () { return exports.LinkRegistry; },
      choiceName: 'CreateTransferPreapproval',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CreateTransferPreapproval.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CreateTransferPreapproval.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkTransferPreapproval).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkTransferPreapproval).encode(__typed__); },
    },
  },
  pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferFactory,
  pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory,
);

damlTypes.registerTemplate(exports.LinkRegistry, ['fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60', '#link']);

exports.LinkTransferInstruction = damlTypes.assembleTemplate(
  {
    templateId: '#link:Link.Token:LinkTransferInstruction',
    templateIdWithPackageId: 'fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60:Link.Token:LinkTransferInstruction',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        instructionAdmin: damlTypes.Party.decoder,
        instructionTransfer: pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.Transfer.decoder,
        instructionLockedHoldingCid: damlTypes.ContractId(exports.LockedLinkHolding).decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        instructionAdmin: damlTypes.Party.encode(__typed__.instructionAdmin),
        instructionTransfer: pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.Transfer.encode(__typed__.instructionTransfer),
        instructionLockedHoldingCid: damlTypes.ContractId(exports.LockedLinkHolding).encode(__typed__.instructionLockedHoldingCid),
      };
    },
    Archive: {
      template: function () { return exports.LinkTransferInstruction; },
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
  },
  pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstruction,
);

damlTypes.registerTemplate(exports.LinkTransferInstruction, ['fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60', '#link']);

exports.LinkTransferPreapproval = damlTypes.assembleTemplate(
  {
    templateId: '#link:Link.Token:LinkTransferPreapproval',
    templateIdWithPackageId: 'fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60:Link.Token:LinkTransferPreapproval',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        preapprovalAdmin: damlTypes.Party.decoder,
        preapprovalReceiver: damlTypes.Party.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        preapprovalAdmin: damlTypes.Party.encode(__typed__.preapprovalAdmin),
        preapprovalReceiver: damlTypes.Party.encode(__typed__.preapprovalReceiver),
      };
    },
    Archive: {
      template: function () { return exports.LinkTransferPreapproval; },
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
    Cancel: {
      template: function () { return exports.LinkTransferPreapproval; },
      choiceName: 'Cancel',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.Cancel.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.Cancel.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    Send: {
      template: function () { return exports.LinkTransferPreapproval; },
      choiceName: 'Send',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.Send.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.Send.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult.decoder;
      }),
      resultEncode: function (__typed__) { return pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult.encode(__typed__); },
    },
  },
);

damlTypes.registerTemplate(exports.LinkTransferPreapproval, ['fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60', '#link']);

exports.LockedLinkHolding = damlTypes.assembleTemplate(
  {
    templateId: '#link:Link.Token:LockedLinkHolding',
    templateIdWithPackageId: 'fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60:Link.Token:LockedLinkHolding',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        lockedAdmin: damlTypes.Party.decoder,
        lockedSender: damlTypes.Party.decoder,
        lockedReceiver: damlTypes.Party.decoder,
        lockedInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        lockedAmount: damlTypes.Numeric(10).decoder,
        meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        lockedAdmin: damlTypes.Party.encode(__typed__.lockedAdmin),
        lockedSender: damlTypes.Party.encode(__typed__.lockedSender),
        lockedReceiver: damlTypes.Party.encode(__typed__.lockedReceiver),
        lockedInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.lockedInstrumentId),
        lockedAmount: damlTypes.Numeric(10).encode(__typed__.lockedAmount),
        meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
      };
    },
    Archive: {
      template: function () { return exports.LockedLinkHolding; },
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
  },
);

damlTypes.registerTemplate(exports.LockedLinkHolding, ['fc1fcb7dbfd1959cafbd54bfa6f6cc88c1678607b0f95ee3f6ee8e9ee8eb2b60', '#link']);

exports.Send = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      sender: damlTypes.Party.decoder,
      amount: damlTypes.Numeric(10).decoder,
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).decoder,
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      sender: damlTypes.Party.encode(__typed__.sender),
      amount: damlTypes.Numeric(10).encode(__typed__.amount),
      instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).encode(__typed__.inputHoldingCids),
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
    };
  },
};
