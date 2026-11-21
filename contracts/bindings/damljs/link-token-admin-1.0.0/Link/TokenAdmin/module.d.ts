// Generated from ../../Link/TokenAdmin/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520 from '@chainlink/canton-contracts/splice-api-token-allocation-instruction-v1-1.0.0';
import * as pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f from '@chainlink/canton-contracts/splice-api-token-metadata-v1-1.0.0';
import * as pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281 from '@chainlink/canton-contracts/splice-api-token-transfer-instruction-v1-1.0.0';
import * as pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240 from '@chainlink/canton-contracts/mcms-api-1.0.0';
import * as pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b from '@chainlink/canton-contracts/splice-api-token-holding-v1-1.0.0';
import * as pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e from '@chainlink/canton-contracts/splice-api-token-burn-mint-v1-1.0.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';

export declare type ApproveAllocate = {
  authorizationId: string,
  receiver: damlTypes.Party,
  amount: damlTypes.Numeric,
  transferLegId: string,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveAllocate:
  damlTypes.Serializable<ApproveAllocate>

export declare type ApproveBurn = {
  authorizationId: string,
  holder: damlTypes.Party,
  maxAmount: damlTypes.Numeric,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveBurn:
  damlTypes.Serializable<ApproveBurn>

export declare type ApproveMint = {
  authorizationId: string,
  recipient: damlTypes.Party,
  amount: damlTypes.Numeric,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveMint:
  damlTypes.Serializable<ApproveMint>

export declare type ApproveTransfer = {
  authorizationId: string,
  receiver: damlTypes.Party,
  amount: damlTypes.Numeric,
  reference: string,
  validForSeconds: damlTypes.Int,
}

export declare const ApproveTransfer:
  damlTypes.Serializable<ApproveTransfer>

export declare type CleanupExpiredAllocationAuthorization = {
  submitter: damlTypes.Party,
}

export declare const CleanupExpiredAllocationAuthorization:
  damlTypes.Serializable<CleanupExpiredAllocationAuthorization>

export declare type CleanupExpiredBurnAuthorization = {
  submitter: damlTypes.Party,
}

export declare const CleanupExpiredBurnAuthorization:
  damlTypes.Serializable<CleanupExpiredBurnAuthorization>

export declare type CleanupExpiredMintAuthorization = {
  submitter: damlTypes.Party,
}

export declare const CleanupExpiredMintAuthorization:
  damlTypes.Serializable<CleanupExpiredMintAuthorization>

export declare type CleanupExpiredTransferAuthorization = {
  submitter: damlTypes.Party,
}

export declare const CleanupExpiredTransferAuthorization:
  damlTypes.Serializable<CleanupExpiredTransferAuthorization>

export declare type DiscardRecovery = {
}

export declare const DiscardRecovery:
  damlTypes.Serializable<DiscardRecovery>

export declare type ExecuteAllocate = {
  submitter: damlTypes.Party,
  tokenAdminCid: damlTypes.ContractId<LinkTokenAdmin>,
  allocationFactoryCid: damlTypes.ContractId<pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520.Splice.Api.Token.AllocationInstructionV1.AllocationFactory>,
  settlementExecutor: damlTypes.Party,
  settlementRefId: string,
  inputHoldingCids: damlTypes.ContractId<pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding>[],
  allocateBeforeHours: damlTypes.Int,
  settleBeforeHours: damlTypes.Int,
  registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare const ExecuteAllocate:
  damlTypes.Serializable<ExecuteAllocate>

export declare type ExecuteBurn = {
  submitter: damlTypes.Party,
  tokenAdminCid: damlTypes.ContractId<LinkTokenAdmin>,
  burnMintFactoryCid: damlTypes.ContractId<pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory>,
  inputHoldingCids: damlTypes.ContractId<pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding>[],
  registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare const ExecuteBurn:
  damlTypes.Serializable<ExecuteBurn>

export declare type ExecuteMint = {
  submitter: damlTypes.Party,
  tokenAdminCid: damlTypes.ContractId<LinkTokenAdmin>,
  burnMintFactoryCid: damlTypes.ContractId<pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory>,
  registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare const ExecuteMint:
  damlTypes.Serializable<ExecuteMint>

export declare type ExecuteTransfer = {
  submitter: damlTypes.Party,
  tokenAdminCid: damlTypes.ContractId<LinkTokenAdmin>,
  transferFactoryCid: damlTypes.ContractId<pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferFactory>,
  inputHoldingCids: damlTypes.ContractId<pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding>[],
  transferTimeoutHours: damlTypes.Optional<damlTypes.Int>,
  registryContext: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.ChoiceContext,
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare const ExecuteTransfer:
  damlTypes.Serializable<ExecuteTransfer>

export declare type LinkAllocationAuthorization = {
  ccipOwner: damlTypes.Party,
  instanceId: string,
  adminInstanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  receiver: damlTypes.Party,
  amount: damlTypes.Numeric,
  transferLegId: string,
  reference: string,
  expiresAt: damlTypes.Time,
}

export declare interface LinkAllocationAuthorizationInterface {
  Archive: 
    damlTypes.Choice<LinkAllocationAuthorization, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkAllocationAuthorization, undefined>>;
  CleanupExpiredAllocationAuthorization: 
    damlTypes.Choice<LinkAllocationAuthorization, CleanupExpiredAllocationAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkAllocationAuthorization, undefined>>;
  ExecuteAllocate: 
    damlTypes.Choice<LinkAllocationAuthorization, ExecuteAllocate, pkg275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520.Splice.Api.Token.AllocationInstructionV1.AllocationInstructionResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkAllocationAuthorization, undefined>>;
  WithdrawAllocationAuthorization: 
    damlTypes.Choice<LinkAllocationAuthorization, WithdrawAllocationAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkAllocationAuthorization, undefined>>;
}
export declare const LinkAllocationAuthorization:
  damlTypes.Template<LinkAllocationAuthorization, undefined, '#link-token-admin:Link.TokenAdmin:LinkAllocationAuthorization'> &
  damlTypes.ToInterface<LinkAllocationAuthorization, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LinkAllocationAuthorizationInterface

export declare type LinkBurnAuthorization = {
  ccipOwner: damlTypes.Party,
  instanceId: string,
  adminInstanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  holder: damlTypes.Party,
  maxAmount: damlTypes.Numeric,
  reference: string,
  expiresAt: damlTypes.Time,
}

export declare interface LinkBurnAuthorizationInterface {
  Archive: 
    damlTypes.Choice<LinkBurnAuthorization, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkBurnAuthorization, undefined>>;
  CleanupExpiredBurnAuthorization: 
    damlTypes.Choice<LinkBurnAuthorization, CleanupExpiredBurnAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkBurnAuthorization, undefined>>;
  ExecuteBurn: 
    damlTypes.Choice<LinkBurnAuthorization, ExecuteBurn, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkBurnAuthorization, undefined>>;
  WithdrawBurnAuthorization: 
    damlTypes.Choice<LinkBurnAuthorization, WithdrawBurnAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkBurnAuthorization, undefined>>;
}
export declare const LinkBurnAuthorization:
  damlTypes.Template<LinkBurnAuthorization, undefined, '#link-token-admin:Link.TokenAdmin:LinkBurnAuthorization'> &
  damlTypes.ToInterface<LinkBurnAuthorization, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LinkBurnAuthorizationInterface

export declare type LinkMintAuthorization = {
  ccipOwner: damlTypes.Party,
  instanceId: string,
  adminInstanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  recipient: damlTypes.Party,
  amount: damlTypes.Numeric,
  reference: string,
  expiresAt: damlTypes.Time,
}

export declare interface LinkMintAuthorizationInterface {
  Archive: 
    damlTypes.Choice<LinkMintAuthorization, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkMintAuthorization, undefined>>;
  CleanupExpiredMintAuthorization: 
    damlTypes.Choice<LinkMintAuthorization, CleanupExpiredMintAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkMintAuthorization, undefined>>;
  ExecuteMint: 
    damlTypes.Choice<LinkMintAuthorization, ExecuteMint, damlTypes.ContractId<pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding>[], undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkMintAuthorization, undefined>>;
  WithdrawMintAuthorization: 
    damlTypes.Choice<LinkMintAuthorization, WithdrawMintAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkMintAuthorization, undefined>>;
}
export declare const LinkMintAuthorization:
  damlTypes.Template<LinkMintAuthorization, undefined, '#link-token-admin:Link.TokenAdmin:LinkMintAuthorization'> &
  damlTypes.ToInterface<LinkMintAuthorization, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LinkMintAuthorizationInterface

export declare type LinkTokenAdmin = {
  ccipOwner: damlTypes.Party,
  instanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  paused: boolean,
  observers: damlTypes.Party[],
}

export declare interface LinkTokenAdminInterface {
  ApproveAllocate: 
    damlTypes.Choice<LinkTokenAdmin, ApproveAllocate, damlTypes.ContractId<LinkAllocationAuthorization>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
  ApproveBurn: 
    damlTypes.Choice<LinkTokenAdmin, ApproveBurn, damlTypes.ContractId<LinkBurnAuthorization>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
  ApproveMint: 
    damlTypes.Choice<LinkTokenAdmin, ApproveMint, damlTypes.ContractId<LinkMintAuthorization>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
  ApproveTransfer: 
    damlTypes.Choice<LinkTokenAdmin, ApproveTransfer, damlTypes.ContractId<LinkTransferAuthorization>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
  Archive: 
    damlTypes.Choice<LinkTokenAdmin, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
  SetObservers: 
    damlTypes.Choice<LinkTokenAdmin, SetObservers, damlTypes.ContractId<LinkTokenAdmin>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
  SetPaused: 
    damlTypes.Choice<LinkTokenAdmin, SetPaused, damlTypes.ContractId<LinkTokenAdmin>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTokenAdmin, undefined>>;
}
export declare const LinkTokenAdmin:
  damlTypes.Template<LinkTokenAdmin, undefined, '#link-token-admin:Link.TokenAdmin:LinkTokenAdmin'> &
  damlTypes.ToInterface<LinkTokenAdmin, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LinkTokenAdminInterface

export declare type LinkTransferAuthorization = {
  ccipOwner: damlTypes.Party,
  instanceId: string,
  adminInstanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  receiver: damlTypes.Party,
  amount: damlTypes.Numeric,
  reference: string,
  expiresAt: damlTypes.Time,
}

export declare interface LinkTransferAuthorizationInterface {
  Archive: 
    damlTypes.Choice<LinkTransferAuthorization, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferAuthorization, undefined>>;
  CleanupExpiredTransferAuthorization: 
    damlTypes.Choice<LinkTransferAuthorization, CleanupExpiredTransferAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferAuthorization, undefined>>;
  ExecuteTransfer: 
    damlTypes.Choice<LinkTransferAuthorization, ExecuteTransfer, pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferAuthorization, undefined>>;
  WithdrawTransferAuthorization: 
    damlTypes.Choice<LinkTransferAuthorization, WithdrawTransferAuthorization, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferAuthorization, undefined>>;
}
export declare const LinkTransferAuthorization:
  damlTypes.Template<LinkTransferAuthorization, undefined, '#link-token-admin:Link.TokenAdmin:LinkTransferAuthorization'> &
  damlTypes.ToInterface<LinkTransferAuthorization, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LinkTransferAuthorizationInterface

export declare type LinkTransferRecovery = {
  ccipOwner: damlTypes.Party,
  instanceId: string,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  transferInstructionCid: damlTypes.ContractId<pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstruction>,
  executeBefore: damlTypes.Time,
}

export declare interface LinkTransferRecoveryInterface {
  Archive: 
    damlTypes.Choice<LinkTransferRecovery, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferRecovery, undefined>>;
  DiscardRecovery: 
    damlTypes.Choice<LinkTransferRecovery, DiscardRecovery, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferRecovery, undefined>>;
  WithdrawInstruction: 
    damlTypes.Choice<LinkTransferRecovery, WithdrawInstruction, pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferRecovery, undefined>>;
}
export declare const LinkTransferRecovery:
  damlTypes.Template<LinkTransferRecovery, undefined, '#link-token-admin:Link.TokenAdmin:LinkTransferRecovery'> &
  damlTypes.ToInterface<LinkTransferRecovery, pkg674d8f60de56afd32698ae19516260217c73dd9ed082680fa840ede4b7665240.MCMS.MCMSReceiver.MCMSReceiver> &
  LinkTransferRecoveryInterface

export declare type SetObservers = {
  observers: damlTypes.Party[],
}

export declare const SetObservers:
  damlTypes.Serializable<SetObservers>

export declare type SetPaused = {
  paused: boolean,
}

export declare const SetPaused:
  damlTypes.Serializable<SetPaused>

export declare type WithdrawAllocationAuthorization = {
}

export declare const WithdrawAllocationAuthorization:
  damlTypes.Serializable<WithdrawAllocationAuthorization>

export declare type WithdrawBurnAuthorization = {
}

export declare const WithdrawBurnAuthorization:
  damlTypes.Serializable<WithdrawBurnAuthorization>

export declare type WithdrawInstruction = {
}

export declare const WithdrawInstruction:
  damlTypes.Serializable<WithdrawInstruction>

export declare type WithdrawMintAuthorization = {
}

export declare const WithdrawMintAuthorization:
  damlTypes.Serializable<WithdrawMintAuthorization>

export declare type WithdrawTransferAuthorization = {
}

export declare const WithdrawTransferAuthorization:
  damlTypes.Serializable<WithdrawTransferAuthorization>
