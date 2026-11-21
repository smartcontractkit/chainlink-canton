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

exports.ApproveAllocateParams = {
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

exports.ApproveBurnParams = {
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

exports.ApproveMintParams = {
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

exports.ApproveTransferParams = {
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

exports.SetObserversParams = {
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

exports.SetPausedParams = {
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
