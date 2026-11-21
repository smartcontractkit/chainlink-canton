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

var pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520 = require('@chainlink/canton-contracts/splice-api-token-allocation-instruction-v1-1.0.0');
var pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f = require('@chainlink/canton-contracts/splice-api-token-metadata-v1-1.0.0');
var pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281 = require('@chainlink/canton-contracts/splice-api-token-transfer-instruction-v1-1.0.0');
var pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 = require('@chainlink/canton-contracts/mcms-api-1.0.0');
var pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b = require('@chainlink/canton-contracts/splice-api-token-holding-v1-1.0.0');
var pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e = require('@chainlink/canton-contracts/splice-api-token-burn-mint-v1-1.0.0');
var pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 = require('@chainlink/canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0');

exports.ApproveAllocate = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      authorizationId: damlTypes.Text.decoder,
      receiver: damlTypes.Party.decoder,
      amount: damlTypes.Numeric(10).decoder,
      transferLegId: damlTypes.Text.decoder,
      reference: damlTypes.Text.decoder,
      validForSeconds: damlTypes.Int.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      authorizationId: damlTypes.Text.encode(__typed__.authorizationId),
      receiver: damlTypes.Party.encode(__typed__.receiver),
      amount: damlTypes.Numeric(10).encode(__typed__.amount),
      transferLegId: damlTypes.Text.encode(__typed__.transferLegId),
      reference: damlTypes.Text.encode(__typed__.reference),
      validForSeconds: damlTypes.Int.encode(__typed__.validForSeconds),
    };
  },
};

exports.ApproveBurn = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      authorizationId: damlTypes.Text.decoder,
      holder: damlTypes.Party.decoder,
      maxAmount: damlTypes.Numeric(10).decoder,
      reference: damlTypes.Text.decoder,
      validForSeconds: damlTypes.Int.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      authorizationId: damlTypes.Text.encode(__typed__.authorizationId),
      holder: damlTypes.Party.encode(__typed__.holder),
      maxAmount: damlTypes.Numeric(10).encode(__typed__.maxAmount),
      reference: damlTypes.Text.encode(__typed__.reference),
      validForSeconds: damlTypes.Int.encode(__typed__.validForSeconds),
    };
  },
};

exports.ApproveMint = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      authorizationId: damlTypes.Text.decoder,
      recipient: damlTypes.Party.decoder,
      amount: damlTypes.Numeric(10).decoder,
      reference: damlTypes.Text.decoder,
      validForSeconds: damlTypes.Int.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      authorizationId: damlTypes.Text.encode(__typed__.authorizationId),
      recipient: damlTypes.Party.encode(__typed__.recipient),
      amount: damlTypes.Numeric(10).encode(__typed__.amount),
      reference: damlTypes.Text.encode(__typed__.reference),
      validForSeconds: damlTypes.Int.encode(__typed__.validForSeconds),
    };
  },
};

exports.ApproveTransfer = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      authorizationId: damlTypes.Text.decoder,
      receiver: damlTypes.Party.decoder,
      amount: damlTypes.Numeric(10).decoder,
      reference: damlTypes.Text.decoder,
      validForSeconds: damlTypes.Int.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      authorizationId: damlTypes.Text.encode(__typed__.authorizationId),
      receiver: damlTypes.Party.encode(__typed__.receiver),
      amount: damlTypes.Numeric(10).encode(__typed__.amount),
      reference: damlTypes.Text.encode(__typed__.reference),
      validForSeconds: damlTypes.Int.encode(__typed__.validForSeconds),
    };
  },
};

exports.CleanupExpiredAllocationAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
    };
  },
};

exports.CleanupExpiredBurnAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
    };
  },
};

exports.CleanupExpiredMintAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
    };
  },
};

exports.CleanupExpiredTransferAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
    };
  },
};

exports.DiscardRecovery = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.ExecuteAllocate = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).decoder,
      allocationFactoryCid: damlTypes.ContractId(pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520.Splice.Api.Token.AllocationInstructionV1.AllocationFactory).decoder,
      settlementExecutor: damlTypes.Party.decoder,
      settlementRefId: damlTypes.Text.decoder,
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).decoder,
      allocateBeforeHours: damlTypes.Int.decoder,
      settleBeforeHours: damlTypes.Int.decoder,
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).encode(__typed__.tokenAdminCid),
      allocationFactoryCid: damlTypes.ContractId(pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520.Splice.Api.Token.AllocationInstructionV1.AllocationFactory).encode(__typed__.allocationFactoryCid),
      settlementExecutor: damlTypes.Party.encode(__typed__.settlementExecutor),
      settlementRefId: damlTypes.Text.encode(__typed__.settlementRefId),
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).encode(__typed__.inputHoldingCids),
      allocateBeforeHours: damlTypes.Int.encode(__typed__.allocateBeforeHours),
      settleBeforeHours: damlTypes.Int.encode(__typed__.settleBeforeHours),
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.registryContext),
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
    };
  },
};

exports.ExecuteBurn = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).decoder,
      burnMintFactoryCid: damlTypes.ContractId(pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory).decoder,
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).decoder,
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).encode(__typed__.tokenAdminCid),
      burnMintFactoryCid: damlTypes.ContractId(pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory).encode(__typed__.burnMintFactoryCid),
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).encode(__typed__.inputHoldingCids),
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.registryContext),
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
    };
  },
};

exports.ExecuteMint = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).decoder,
      burnMintFactoryCid: damlTypes.ContractId(pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory).decoder,
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).encode(__typed__.tokenAdminCid),
      burnMintFactoryCid: damlTypes.ContractId(pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory).encode(__typed__.burnMintFactoryCid),
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.registryContext),
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
    };
  },
};

exports.ExecuteTransfer = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).decoder,
      transferFactoryCid: damlTypes.ContractId(pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferFactory).decoder,
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).decoder,
      transferTimeoutHours: jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Int).decoder),
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.decoder,
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      tokenAdminCid: damlTypes.ContractId(exports.LinkTokenAdmin).encode(__typed__.tokenAdminCid),
      transferFactoryCid: damlTypes.ContractId(pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferFactory).encode(__typed__.transferFactoryCid),
      inputHoldingCids: damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).encode(__typed__.inputHoldingCids),
      transferTimeoutHours: damlTypes.Optional(damlTypes.Int).encode(__typed__.transferTimeoutHours),
      registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext.encode(__typed__.registryContext),
      meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata.encode(__typed__.meta),
    };
  },
};

exports.LinkAllocationAuthorization = damlTypes.assembleTemplate(
  {
    templateId: '#link-token-admin:Link.TokenAdmin:LinkAllocationAuthorization',
    templateIdWithPackageId: '506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a:Link.TokenAdmin:LinkAllocationAuthorization',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        ccipOwner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        adminInstanceId: damlTypes.Text.decoder,
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        receiver: damlTypes.Party.decoder,
        amount: damlTypes.Numeric(10).decoder,
        transferLegId: damlTypes.Text.decoder,
        reference: damlTypes.Text.decoder,
        expiresAt: damlTypes.Time.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        adminInstanceId: damlTypes.Text.encode(__typed__.adminInstanceId),
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
        receiver: damlTypes.Party.encode(__typed__.receiver),
        amount: damlTypes.Numeric(10).encode(__typed__.amount),
        transferLegId: damlTypes.Text.encode(__typed__.transferLegId),
        reference: damlTypes.Text.encode(__typed__.reference),
        expiresAt: damlTypes.Time.encode(__typed__.expiresAt),
      };
    },
    Archive: {
      template: function () { return exports.LinkAllocationAuthorization; },
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
    CleanupExpiredAllocationAuthorization: {
      template: function () { return exports.LinkAllocationAuthorization; },
      choiceName: 'CleanupExpiredAllocationAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CleanupExpiredAllocationAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CleanupExpiredAllocationAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    ExecuteAllocate: {
      template: function () { return exports.LinkAllocationAuthorization; },
      choiceName: 'ExecuteAllocate',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ExecuteAllocate.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ExecuteAllocate.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520.Splice.Api.Token.AllocationInstructionV1.AllocationInstructionResult.decoder;
      }),
      resultEncode: function (__typed__) { return pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520.Splice.Api.Token.AllocationInstructionV1.AllocationInstructionResult.encode(__typed__); },
    },
    WithdrawAllocationAuthorization: {
      template: function () { return exports.LinkAllocationAuthorization; },
      choiceName: 'WithdrawAllocationAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.WithdrawAllocationAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.WithdrawAllocationAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.LinkAllocationAuthorization, ['506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a', '#link-token-admin']);

exports.LinkBurnAuthorization = damlTypes.assembleTemplate(
  {
    templateId: '#link-token-admin:Link.TokenAdmin:LinkBurnAuthorization',
    templateIdWithPackageId: '506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a:Link.TokenAdmin:LinkBurnAuthorization',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        ccipOwner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        adminInstanceId: damlTypes.Text.decoder,
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        holder: damlTypes.Party.decoder,
        maxAmount: damlTypes.Numeric(10).decoder,
        reference: damlTypes.Text.decoder,
        expiresAt: damlTypes.Time.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        adminInstanceId: damlTypes.Text.encode(__typed__.adminInstanceId),
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
        holder: damlTypes.Party.encode(__typed__.holder),
        maxAmount: damlTypes.Numeric(10).encode(__typed__.maxAmount),
        reference: damlTypes.Text.encode(__typed__.reference),
        expiresAt: damlTypes.Time.encode(__typed__.expiresAt),
      };
    },
    Archive: {
      template: function () { return exports.LinkBurnAuthorization; },
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
    CleanupExpiredBurnAuthorization: {
      template: function () { return exports.LinkBurnAuthorization; },
      choiceName: 'CleanupExpiredBurnAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CleanupExpiredBurnAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CleanupExpiredBurnAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    ExecuteBurn: {
      template: function () { return exports.LinkBurnAuthorization; },
      choiceName: 'ExecuteBurn',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ExecuteBurn.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ExecuteBurn.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    WithdrawBurnAuthorization: {
      template: function () { return exports.LinkBurnAuthorization; },
      choiceName: 'WithdrawBurnAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.WithdrawBurnAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.WithdrawBurnAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.LinkBurnAuthorization, ['506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a', '#link-token-admin']);

exports.LinkMintAuthorization = damlTypes.assembleTemplate(
  {
    templateId: '#link-token-admin:Link.TokenAdmin:LinkMintAuthorization',
    templateIdWithPackageId: '506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a:Link.TokenAdmin:LinkMintAuthorization',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        ccipOwner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        adminInstanceId: damlTypes.Text.decoder,
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        recipient: damlTypes.Party.decoder,
        amount: damlTypes.Numeric(10).decoder,
        reference: damlTypes.Text.decoder,
        expiresAt: damlTypes.Time.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        adminInstanceId: damlTypes.Text.encode(__typed__.adminInstanceId),
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
        recipient: damlTypes.Party.encode(__typed__.recipient),
        amount: damlTypes.Numeric(10).encode(__typed__.amount),
        reference: damlTypes.Text.encode(__typed__.reference),
        expiresAt: damlTypes.Time.encode(__typed__.expiresAt),
      };
    },
    Archive: {
      template: function () { return exports.LinkMintAuthorization; },
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
    CleanupExpiredMintAuthorization: {
      template: function () { return exports.LinkMintAuthorization; },
      choiceName: 'CleanupExpiredMintAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CleanupExpiredMintAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CleanupExpiredMintAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    ExecuteMint: {
      template: function () { return exports.LinkMintAuthorization; },
      choiceName: 'ExecuteMint',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ExecuteMint.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ExecuteMint.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.List(damlTypes.ContractId(pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding)).encode(__typed__); },
    },
    WithdrawMintAuthorization: {
      template: function () { return exports.LinkMintAuthorization; },
      choiceName: 'WithdrawMintAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.WithdrawMintAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.WithdrawMintAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.LinkMintAuthorization, ['506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a', '#link-token-admin']);

exports.LinkTokenAdmin = damlTypes.assembleTemplate(
  {
    templateId: '#link-token-admin:Link.TokenAdmin:LinkTokenAdmin',
    templateIdWithPackageId: '506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a:Link.TokenAdmin:LinkTokenAdmin',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        ccipOwner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        paused: damlTypes.Bool.decoder,
        observers: damlTypes.List(damlTypes.Party).decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
        paused: damlTypes.Bool.encode(__typed__.paused),
        observers: damlTypes.List(damlTypes.Party).encode(__typed__.observers),
      };
    },
    ApproveAllocate: {
      template: function () { return exports.LinkTokenAdmin; },
      choiceName: 'ApproveAllocate',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApproveAllocate.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApproveAllocate.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkAllocationAuthorization).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkAllocationAuthorization).encode(__typed__); },
    },
    ApproveBurn: {
      template: function () { return exports.LinkTokenAdmin; },
      choiceName: 'ApproveBurn',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApproveBurn.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApproveBurn.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkBurnAuthorization).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkBurnAuthorization).encode(__typed__); },
    },
    ApproveMint: {
      template: function () { return exports.LinkTokenAdmin; },
      choiceName: 'ApproveMint',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApproveMint.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApproveMint.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkMintAuthorization).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkMintAuthorization).encode(__typed__); },
    },
    ApproveTransfer: {
      template: function () { return exports.LinkTokenAdmin; },
      choiceName: 'ApproveTransfer',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ApproveTransfer.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ApproveTransfer.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkTransferAuthorization).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkTransferAuthorization).encode(__typed__); },
    },
    Archive: {
      template: function () { return exports.LinkTokenAdmin; },
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
    SetObservers: {
      template: function () { return exports.LinkTokenAdmin; },
      choiceName: 'SetObservers',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetObservers.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetObservers.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkTokenAdmin).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkTokenAdmin).encode(__typed__); },
    },
    SetPaused: {
      template: function () { return exports.LinkTokenAdmin; },
      choiceName: 'SetPaused',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetPaused.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetPaused.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.LinkTokenAdmin).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.LinkTokenAdmin).encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.LinkTokenAdmin, ['506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a', '#link-token-admin']);

exports.LinkTransferAuthorization = damlTypes.assembleTemplate(
  {
    templateId: '#link-token-admin:Link.TokenAdmin:LinkTransferAuthorization',
    templateIdWithPackageId: '506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a:Link.TokenAdmin:LinkTransferAuthorization',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        ccipOwner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        adminInstanceId: damlTypes.Text.decoder,
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        receiver: damlTypes.Party.decoder,
        amount: damlTypes.Numeric(10).decoder,
        reference: damlTypes.Text.decoder,
        expiresAt: damlTypes.Time.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        adminInstanceId: damlTypes.Text.encode(__typed__.adminInstanceId),
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
        receiver: damlTypes.Party.encode(__typed__.receiver),
        amount: damlTypes.Numeric(10).encode(__typed__.amount),
        reference: damlTypes.Text.encode(__typed__.reference),
        expiresAt: damlTypes.Time.encode(__typed__.expiresAt),
      };
    },
    Archive: {
      template: function () { return exports.LinkTransferAuthorization; },
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
    CleanupExpiredTransferAuthorization: {
      template: function () { return exports.LinkTransferAuthorization; },
      choiceName: 'CleanupExpiredTransferAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CleanupExpiredTransferAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CleanupExpiredTransferAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    ExecuteTransfer: {
      template: function () { return exports.LinkTransferAuthorization; },
      choiceName: 'ExecuteTransfer',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ExecuteTransfer.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ExecuteTransfer.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult.decoder;
      }),
      resultEncode: function (__typed__) { return pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult.encode(__typed__); },
    },
    WithdrawTransferAuthorization: {
      template: function () { return exports.LinkTransferAuthorization; },
      choiceName: 'WithdrawTransferAuthorization',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.WithdrawTransferAuthorization.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.WithdrawTransferAuthorization.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.LinkTransferAuthorization, ['506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a', '#link-token-admin']);

exports.LinkTransferRecovery = damlTypes.assembleTemplate(
  {
    templateId: '#link-token-admin:Link.TokenAdmin:LinkTransferRecovery',
    templateIdWithPackageId: '506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a:Link.TokenAdmin:LinkTransferRecovery',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        ccipOwner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.decoder,
        transferInstructionCid: damlTypes.ContractId(pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstruction).decoder,
        executeBefore: damlTypes.Time.decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        ccipOwner: damlTypes.Party.encode(__typed__.ccipOwner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId.encode(__typed__.instrumentId),
        transferInstructionCid: damlTypes.ContractId(pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstruction).encode(__typed__.transferInstructionCid),
        executeBefore: damlTypes.Time.encode(__typed__.executeBefore),
      };
    },
    Archive: {
      template: function () { return exports.LinkTransferRecovery; },
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
    DiscardRecovery: {
      template: function () { return exports.LinkTransferRecovery; },
      choiceName: 'DiscardRecovery',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.DiscardRecovery.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.DiscardRecovery.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Unit.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Unit.encode(__typed__); },
    },
    WithdrawInstruction: {
      template: function () { return exports.LinkTransferRecovery; },
      choiceName: 'WithdrawInstruction',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.WithdrawInstruction.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.WithdrawInstruction.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult.decoder;
      }),
      resultEncode: function (__typed__) { return pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult.encode(__typed__); },
    },
  },
  pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver,
);

damlTypes.registerTemplate(exports.LinkTransferRecovery, ['506338cda467b8dc4515e66cddaabd234e8d3f9bdb4a3222f252e205c7afd66a', '#link-token-admin']);

exports.SetObservers = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      observers: damlTypes.List(damlTypes.Party).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      observers: damlTypes.List(damlTypes.Party).encode(__typed__.observers),
    };
  },
};

exports.SetPaused = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      paused: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      paused: damlTypes.Bool.encode(__typed__.paused),
    };
  },
};

exports.WithdrawAllocationAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.WithdrawBurnAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.WithdrawInstruction = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.WithdrawMintAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.WithdrawTransferAuthorization = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};
