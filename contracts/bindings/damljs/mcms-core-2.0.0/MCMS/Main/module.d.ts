// Generated from ../../MCMS/Main/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 from '@chainlink/canton-contracts/mcms-api-1.0.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';
import * as pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946 from '@chainlink/canton-contracts/daml-stdlib-DA-Time-Types-1.0.0';

export declare type CanExecuteOp = {
  submitter: damlTypes.Party,
  targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
  op: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Op,
}

export declare const CanExecuteOp:
  damlTypes.Serializable<CanExecuteOp>

export declare type ExecuteOp = {
  targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
  submitter: damlTypes.Party,
  op: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Op,
  opProof: string[],
  targetCids: damlTypes.Map<string, damlTypes.ContractId<pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver>>,
}

export declare const ExecuteOp:
  damlTypes.Serializable<ExecuteOp>

export declare type ExecuteScheduledBatch = {
  submitter: damlTypes.Party,
  opId: string,
  calls: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.TimelockCall[],
  predecessor: string,
  salt: string,
  targetCids: damlTypes.Map<string, damlTypes.ContractId<pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver>>,
}

export declare const ExecuteScheduledBatch:
  damlTypes.Serializable<ExecuteScheduledBatch>

export declare type GetBlockedFunctions = {
  submitter: damlTypes.Party,
}

export declare const GetBlockedFunctions:
  damlTypes.Serializable<GetBlockedFunctions>

export declare type GetBlockedFunctionsCount = {
  submitter: damlTypes.Party,
}

export declare const GetBlockedFunctionsCount:
  damlTypes.Serializable<GetBlockedFunctionsCount>

export declare type GetMinDelay = {
  submitter: damlTypes.Party,
}

export declare const GetMinDelay:
  damlTypes.Serializable<GetMinDelay>

export declare type GetState = {
  submitter: damlTypes.Party,
  targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
}

export declare const GetState:
  damlTypes.Serializable<GetState>

export declare type GetTimestamp = {
  submitter: damlTypes.Party,
  opId: string,
}

export declare const GetTimestamp:
  damlTypes.Serializable<GetTimestamp>

export declare type IsOperation = {
  submitter: damlTypes.Party,
  opId: string,
}

export declare const IsOperation:
  damlTypes.Serializable<IsOperation>

export declare type IsOperationDone = {
  submitter: damlTypes.Party,
  opId: string,
}

export declare const IsOperationDone:
  damlTypes.Serializable<IsOperationDone>

export declare type IsOperationPending = {
  submitter: damlTypes.Party,
  opId: string,
}

export declare const IsOperationPending:
  damlTypes.Serializable<IsOperationPending>

export declare type IsOperationReady = {
  submitter: damlTypes.Party,
  opId: string,
}

export declare const IsOperationReady:
  damlTypes.Serializable<IsOperationReady>

export declare type MCMS = {
  owner: damlTypes.Party,
  instanceId: string,
  chainId: damlTypes.Int,
  proposer: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState,
  canceller: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState,
  bypasser: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RoleState,
  minDelay: pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946.DA.Time.Types.RelTime,
  blockedFunctions: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.BlockedFunction[],
  timelockTimestamps: damlTypes.Map<string, damlTypes.Time>,
}

export declare interface MCMSInterface {
  Archive: 
    damlTypes.Choice<MCMS, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  CanExecuteOp: 
    damlTypes.Choice<MCMS, CanExecuteOp, boolean, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  ExecuteOp: 
    damlTypes.Choice<MCMS, ExecuteOp, damlTypes.ContractId<MCMS>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  ExecuteScheduledBatch: 
    damlTypes.Choice<MCMS, ExecuteScheduledBatch, damlTypes.ContractId<MCMS>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  GetBlockedFunctions: 
    damlTypes.Choice<MCMS, GetBlockedFunctions, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.BlockedFunction[], undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  GetBlockedFunctionsCount: 
    damlTypes.Choice<MCMS, GetBlockedFunctionsCount, damlTypes.Int, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  GetMinDelay: 
    damlTypes.Choice<MCMS, GetMinDelay, pkgb70db8369e1c461d5c70f1c86f526a29e9776c655e6ffc2560f95b05ccb8b946.DA.Time.Types.RelTime, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  GetState: 
    damlTypes.Choice<MCMS, GetState, MCMSState, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  GetTimestamp: 
    damlTypes.Choice<MCMS, GetTimestamp, damlTypes.Optional<damlTypes.Time>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  IsOperation: 
    damlTypes.Choice<MCMS, IsOperation, boolean, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  IsOperationDone: 
    damlTypes.Choice<MCMS, IsOperationDone, boolean, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  IsOperationPending: 
    damlTypes.Choice<MCMS, IsOperationPending, boolean, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  IsOperationReady: 
    damlTypes.Choice<MCMS, IsOperationReady, boolean, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  PruneSeenHashes: 
    damlTypes.Choice<MCMS, PruneSeenHashes, damlTypes.ContractId<MCMS>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  PruneTimelockTimestamps: 
    damlTypes.Choice<MCMS, PruneTimelockTimestamps, damlTypes.ContractId<MCMS>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  SetConfig: 
    damlTypes.Choice<MCMS, SetConfig, damlTypes.ContractId<MCMS>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
  SetRoot: 
    damlTypes.Choice<MCMS, SetRoot, damlTypes.ContractId<MCMS>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<MCMS, undefined>>;
}
export declare const MCMS:
  damlTypes.Template<MCMS, undefined, '#mcms-core:MCMS.Main:MCMS'> &
  damlTypes.ToInterface<MCMS, never> &
  MCMSInterface

export declare type MCMSState = {
  role: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
  opCount: damlTypes.Int,
  postOpCount: damlTypes.Int,
  validUntil: damlTypes.Time,
  hasActiveRoot: boolean,
  numSigners: damlTypes.Int,
}

export declare const MCMSState:
  damlTypes.Serializable<MCMSState>

export declare type PruneSeenHashes = {
  targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
}

export declare const PruneSeenHashes:
  damlTypes.Serializable<PruneSeenHashes>

export declare type PruneTimelockTimestamps = {
}

export declare const PruneTimelockTimestamps:
  damlTypes.Serializable<PruneTimelockTimestamps>

export declare type SetConfig = {
  targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
  newSigners: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.SignerInfo[],
  newGroupQuorums: damlTypes.Int[],
  newGroupParents: damlTypes.Int[],
  clearRoot: boolean,
}

export declare const SetConfig:
  damlTypes.Serializable<SetConfig>

export declare type SetRoot = {
  targetRole: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.Role,
  submitter: damlTypes.Party,
  newRoot: string,
  validUntil: damlTypes.Time,
  metadata: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Types.RootMetadata,
  metadataProof: string[],
  signatures: pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.Crypto.RawSignature[],
}

export declare const SetRoot:
  damlTypes.Serializable<SetRoot>
