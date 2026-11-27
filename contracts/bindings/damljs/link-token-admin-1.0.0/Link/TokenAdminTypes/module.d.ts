// Generated from ../../Link/TokenAdminTypes/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

export declare type ApproveAllocateParams = {
  authorizationId: string,
  receiver: damlTypes.Party,
  amount: damlTypes.Numeric,
  transferLegId: string,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveAllocateParams:
  damlTypes.Serializable<ApproveAllocateParams>

export declare type ApproveBurnParams = {
  authorizationId: string,
  holder: damlTypes.Party,
  maxAmount: damlTypes.Numeric,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveBurnParams:
  damlTypes.Serializable<ApproveBurnParams>

export declare type ApproveMintParams = {
  authorizationId: string,
  recipient: damlTypes.Party,
  amount: damlTypes.Numeric,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveMintParams:
  damlTypes.Serializable<ApproveMintParams>

export declare type ApproveTransferParams = {
  authorizationId: string,
  receiver: damlTypes.Party,
  amount: damlTypes.Numeric,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveTransferParams:
  damlTypes.Serializable<ApproveTransferParams>

export declare type SetObserversParams = {
  observers: damlTypes.Party[],
}

export declare const SetObserversParams:
  damlTypes.Serializable<SetObserversParams>

export declare type SetPausedParams = {
  paused: boolean,
}

export declare const SetPausedParams:
  damlTypes.Serializable<SetPausedParams>
