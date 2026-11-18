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

var pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 = require('@chainlink/chainlink-canton-contracts/mcms-api-1.0.0');
var pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 = require('@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0');
var pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946 = require('@chainlink/chainlink-canton-contracts/daml-stdlib-DA-Time-Types-1.0.0');

exports.CanExecuteOp = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
      op: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Op.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.targetRole),
      op: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Op.encode(__typed__.op),
    };
  },
};

exports.ExecuteOp = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
      submitter: damlTypes.Party.decoder,
      op: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Op.decoder,
      opProof: damlTypes.List(damlTypes.Text).decoder,
      targetCids: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver)).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.targetRole),
      submitter: damlTypes.Party.encode(__typed__.submitter),
      op: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Op.encode(__typed__.op),
      opProof: damlTypes.List(damlTypes.Text).encode(__typed__.opProof),
      targetCids: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver)).encode(__typed__.targetCids),
    };
  },
};

exports.ExecuteScheduledBatch = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      opId: damlTypes.Text.decoder,
      calls: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.TimelockCall).decoder,
      predecessor: damlTypes.Text.decoder,
      salt: damlTypes.Text.decoder,
      targetCids: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver)).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      opId: damlTypes.Text.encode(__typed__.opId),
      calls: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.TimelockCall).encode(__typed__.calls),
      predecessor: damlTypes.Text.encode(__typed__.predecessor),
      salt: damlTypes.Text.encode(__typed__.salt),
      targetCids: damlTypes.Map(damlTypes.Text, damlTypes.ContractId(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver)).encode(__typed__.targetCids),
    };
  },
};

exports.GetBlockedFunctions = {
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

exports.GetBlockedFunctionsCount = {
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

exports.GetMinDelay = {
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

exports.GetState = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.targetRole),
    };
  },
};

exports.GetTimestamp = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      opId: damlTypes.Text.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      opId: damlTypes.Text.encode(__typed__.opId),
    };
  },
};

exports.IsOperation = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      opId: damlTypes.Text.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      opId: damlTypes.Text.encode(__typed__.opId),
    };
  },
};

exports.IsOperationDone = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      opId: damlTypes.Text.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      opId: damlTypes.Text.encode(__typed__.opId),
    };
  },
};

exports.IsOperationPending = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      opId: damlTypes.Text.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      opId: damlTypes.Text.encode(__typed__.opId),
    };
  },
};

exports.IsOperationReady = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      submitter: damlTypes.Party.decoder,
      opId: damlTypes.Text.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      submitter: damlTypes.Party.encode(__typed__.submitter),
      opId: damlTypes.Text.encode(__typed__.opId),
    };
  },
};

exports.MCMS = damlTypes.assembleTemplate(
  {
    templateId: '#mcms-core:MCMS.Main:MCMS',
    templateIdWithPackageId: '001aaaed38e8f24dca44c618ccfe7779e9ad2eb120395c89b7d9828617144f5c:MCMS.Main:MCMS',
    keyDecoder: jtv.constant(undefined),
    keyEncode: function () { throw 'EncodeError'; },
    decoder: damlTypes.lazyMemo(function () {
      return jtv.object({
        owner: damlTypes.Party.decoder,
        instanceId: damlTypes.Text.decoder,
        chainId: damlTypes.Int.decoder,
        proposer: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState.decoder,
        canceller: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState.decoder,
        bypasser: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState.decoder,
        minDelay: pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946.DA.Time.Types.RelTime.decoder,
        blockedFunctions: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.BlockedFunction).decoder,
        timelockTimestamps: damlTypes.Map(damlTypes.Text, damlTypes.Time).decoder,
      });
    }),
    encode: function (__typed__) {
      return {
        owner: damlTypes.Party.encode(__typed__.owner),
        instanceId: damlTypes.Text.encode(__typed__.instanceId),
        chainId: damlTypes.Int.encode(__typed__.chainId),
        proposer: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState.encode(__typed__.proposer),
        canceller: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState.encode(__typed__.canceller),
        bypasser: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState.encode(__typed__.bypasser),
        minDelay: pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946.DA.Time.Types.RelTime.encode(__typed__.minDelay),
        blockedFunctions: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.BlockedFunction).encode(__typed__.blockedFunctions),
        timelockTimestamps: damlTypes.Map(damlTypes.Text, damlTypes.Time).encode(__typed__.timelockTimestamps),
      };
    },
    Archive: {
      template: function () { return exports.MCMS; },
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
    CanExecuteOp: {
      template: function () { return exports.MCMS; },
      choiceName: 'CanExecuteOp',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.CanExecuteOp.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.CanExecuteOp.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Bool.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Bool.encode(__typed__); },
    },
    ExecuteOp: {
      template: function () { return exports.MCMS; },
      choiceName: 'ExecuteOp',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ExecuteOp.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ExecuteOp.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.MCMS).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.MCMS).encode(__typed__); },
    },
    ExecuteScheduledBatch: {
      template: function () { return exports.MCMS; },
      choiceName: 'ExecuteScheduledBatch',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.ExecuteScheduledBatch.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.ExecuteScheduledBatch.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.MCMS).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.MCMS).encode(__typed__); },
    },
    GetBlockedFunctions: {
      template: function () { return exports.MCMS; },
      choiceName: 'GetBlockedFunctions',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetBlockedFunctions.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetBlockedFunctions.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.BlockedFunction).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.BlockedFunction).encode(__typed__); },
    },
    GetBlockedFunctionsCount: {
      template: function () { return exports.MCMS; },
      choiceName: 'GetBlockedFunctionsCount',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetBlockedFunctionsCount.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetBlockedFunctionsCount.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Int.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Int.encode(__typed__); },
    },
    GetMinDelay: {
      template: function () { return exports.MCMS; },
      choiceName: 'GetMinDelay',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetMinDelay.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetMinDelay.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946.DA.Time.Types.RelTime.decoder;
      }),
      resultEncode: function (__typed__) { return pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946.DA.Time.Types.RelTime.encode(__typed__); },
    },
    GetState: {
      template: function () { return exports.MCMS; },
      choiceName: 'GetState',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetState.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetState.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return exports.MCMSState.decoder;
      }),
      resultEncode: function (__typed__) { return exports.MCMSState.encode(__typed__); },
    },
    GetTimestamp: {
      template: function () { return exports.MCMS; },
      choiceName: 'GetTimestamp',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.GetTimestamp.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.GetTimestamp.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return jtv.Decoder.withDefault(null, damlTypes.Optional(damlTypes.Time).decoder);
      }),
      resultEncode: function (__typed__) { return damlTypes.Optional(damlTypes.Time).encode(__typed__); },
    },
    IsOperation: {
      template: function () { return exports.MCMS; },
      choiceName: 'IsOperation',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.IsOperation.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.IsOperation.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Bool.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Bool.encode(__typed__); },
    },
    IsOperationDone: {
      template: function () { return exports.MCMS; },
      choiceName: 'IsOperationDone',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.IsOperationDone.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.IsOperationDone.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Bool.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Bool.encode(__typed__); },
    },
    IsOperationPending: {
      template: function () { return exports.MCMS; },
      choiceName: 'IsOperationPending',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.IsOperationPending.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.IsOperationPending.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Bool.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Bool.encode(__typed__); },
    },
    IsOperationReady: {
      template: function () { return exports.MCMS; },
      choiceName: 'IsOperationReady',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.IsOperationReady.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.IsOperationReady.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.Bool.decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.Bool.encode(__typed__); },
    },
    PruneSeenHashes: {
      template: function () { return exports.MCMS; },
      choiceName: 'PruneSeenHashes',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.PruneSeenHashes.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.PruneSeenHashes.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.MCMS).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.MCMS).encode(__typed__); },
    },
    PruneTimelockTimestamps: {
      template: function () { return exports.MCMS; },
      choiceName: 'PruneTimelockTimestamps',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.PruneTimelockTimestamps.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.PruneTimelockTimestamps.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.MCMS).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.MCMS).encode(__typed__); },
    },
    SetConfig: {
      template: function () { return exports.MCMS; },
      choiceName: 'SetConfig',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetConfig.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetConfig.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.MCMS).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.MCMS).encode(__typed__); },
    },
    SetRoot: {
      template: function () { return exports.MCMS; },
      choiceName: 'SetRoot',
      argumentDecoder: damlTypes.lazyMemo(function () {
        return exports.SetRoot.decoder;
      }),
      argumentEncode: function (__typed__) { return exports.SetRoot.encode(__typed__); },
      resultDecoder: damlTypes.lazyMemo(function () {
        return damlTypes.ContractId(exports.MCMS).decoder;
      }),
      resultEncode: function (__typed__) { return damlTypes.ContractId(exports.MCMS).encode(__typed__); },
    },
  },
);

damlTypes.registerTemplate(exports.MCMS, ['001aaaed38e8f24dca44c618ccfe7779e9ad2eb120395c89b7d9828617144f5c', '#mcms-core']);

exports.MCMSState = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      role: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
      opCount: damlTypes.Int.decoder,
      postOpCount: damlTypes.Int.decoder,
      validUntil: damlTypes.Time.decoder,
      hasActiveRoot: damlTypes.Bool.decoder,
      numSigners: damlTypes.Int.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      role: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.role),
      opCount: damlTypes.Int.encode(__typed__.opCount),
      postOpCount: damlTypes.Int.encode(__typed__.postOpCount),
      validUntil: damlTypes.Time.encode(__typed__.validUntil),
      hasActiveRoot: damlTypes.Bool.encode(__typed__.hasActiveRoot),
      numSigners: damlTypes.Int.encode(__typed__.numSigners),
    };
  },
};

exports.PruneSeenHashes = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.targetRole),
    };
  },
};

exports.PruneTimelockTimestamps = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
    });
  }),
  encode: function (__typed__) {
    return {};
  },
};

exports.SetConfig = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
      newSigners: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.SignerInfo).decoder,
      newGroupQuorums: damlTypes.List(damlTypes.Int).decoder,
      newGroupParents: damlTypes.List(damlTypes.Int).decoder,
      clearRoot: damlTypes.Bool.decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.targetRole),
      newSigners: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.SignerInfo).encode(__typed__.newSigners),
      newGroupQuorums: damlTypes.List(damlTypes.Int).encode(__typed__.newGroupQuorums),
      newGroupParents: damlTypes.List(damlTypes.Int).encode(__typed__.newGroupParents),
      clearRoot: damlTypes.Bool.encode(__typed__.clearRoot),
    };
  },
};

exports.SetRoot = {
  decoder: damlTypes.lazyMemo(function () {
    return jtv.object({
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.decoder,
      submitter: damlTypes.Party.decoder,
      newRoot: damlTypes.Text.decoder,
      validUntil: damlTypes.Time.decoder,
      metadata: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RootMetadata.decoder,
      metadataProof: damlTypes.List(damlTypes.Text).decoder,
      signatures: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Crypto.RawSignature).decoder,
    });
  }),
  encode: function (__typed__) {
    return {
      targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role.encode(__typed__.targetRole),
      submitter: damlTypes.Party.encode(__typed__.submitter),
      newRoot: damlTypes.Text.encode(__typed__.newRoot),
      validUntil: damlTypes.Time.encode(__typed__.validUntil),
      metadata: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RootMetadata.encode(__typed__.metadata),
      metadataProof: damlTypes.List(damlTypes.Text).encode(__typed__.metadataProof),
      signatures: damlTypes.List(pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Crypto.RawSignature).encode(__typed__.signatures),
    };
  },
};
