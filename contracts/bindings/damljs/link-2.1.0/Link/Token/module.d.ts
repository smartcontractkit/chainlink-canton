// Generated from ../../Link/Token/module.daml

/* eslint-disable @typescript-eslint/camelcase */
/* eslint-disable @typescript-eslint/no-namespace */
/* eslint-disable @typescript-eslint/no-use-before-define */
import * as jtv from '@mojotech/json-type-validation';
import * as damlTypes from '@daml/types';

import * as pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f from '@chainlink/chainlink-canton-contracts/splice-api-token-metadata-v1-1.0.0';
import * as pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281 from '@chainlink/chainlink-canton-contracts/splice-api-token-transfer-instruction-v1-1.0.0';
import * as pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b from '@chainlink/chainlink-canton-contracts/splice-api-token-holding-v1-1.0.0';
import * as pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e from '@chainlink/chainlink-canton-contracts/splice-api-token-burn-mint-v1-1.0.0';
import * as pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69 from '@chainlink/chainlink-canton-contracts/ghc-stdlib-DA-Internal-Template-1.0.0';

export declare type Cancel = {
  actor: damlTypes.Party,
}

export declare const Cancel:
  damlTypes.Serializable<Cancel>

export declare type CreateTransferPreapproval = {
  receiver: damlTypes.Party,
  existingPreapprovalCid: damlTypes.Optional<damlTypes.ContractId<LinkTransferPreapproval>>,
}

export declare const CreateTransferPreapproval:
  damlTypes.Serializable<CreateTransferPreapproval>

export declare type LinkHolding = {
  holdingOwner: damlTypes.Party,
  holdingAdmin: damlTypes.Party,
  holdingInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  holdingAmount: damlTypes.Numeric,
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare interface LinkHoldingInterface {
  Archive: 
    damlTypes.Choice<LinkHolding, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkHolding, undefined>>;
}
export declare const LinkHolding:
  damlTypes.Template<LinkHolding, undefined, '#link:Link.Token:LinkHolding'> &
  damlTypes.ToInterface<LinkHolding, pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding> &
  LinkHoldingInterface

export declare type LinkRegistry = {
  registryAdmin: damlTypes.Party,
  registryInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  instanceId: string,
  registryMeta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare interface LinkRegistryInterface {
  Archive: 
    damlTypes.Choice<LinkRegistry, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkRegistry, undefined>>;
  CreateTransferPreapproval: 
    damlTypes.Choice<LinkRegistry, CreateTransferPreapproval, damlTypes.ContractId<LinkTransferPreapproval>, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkRegistry, undefined>>;
}
export declare const LinkRegistry:
  damlTypes.Template<LinkRegistry, undefined, '#link:Link.Token:LinkRegistry'> &
  damlTypes.ToInterface<LinkRegistry, pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferFactory | pkg9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e.Splice.Api.Token.BurnMintV1.BurnMintFactory> &
  LinkRegistryInterface

export declare type LinkTransferInstruction = {
  instructionAdmin: damlTypes.Party,
  instructionTransfer: pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.Transfer,
  instructionLockedHoldingCid: damlTypes.ContractId<LockedLinkHolding>,
}

export declare interface LinkTransferInstructionInterface {
  Archive: 
    damlTypes.Choice<LinkTransferInstruction, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferInstruction, undefined>>;
}
export declare const LinkTransferInstruction:
  damlTypes.Template<LinkTransferInstruction, undefined, '#link:Link.Token:LinkTransferInstruction'> &
  damlTypes.ToInterface<LinkTransferInstruction, pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstruction> &
  LinkTransferInstructionInterface

export declare type LinkTransferPreapproval = {
  preapprovalAdmin: damlTypes.Party,
  preapprovalReceiver: damlTypes.Party,
}

export declare interface LinkTransferPreapprovalInterface {
  Archive: 
    damlTypes.Choice<LinkTransferPreapproval, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferPreapproval, undefined>>;
  Cancel: 
    damlTypes.Choice<LinkTransferPreapproval, Cancel, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferPreapproval, undefined>>;
  Send: 
    damlTypes.Choice<LinkTransferPreapproval, Send, pkg55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281.Splice.Api.Token.TransferInstructionV1.TransferInstructionResult, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LinkTransferPreapproval, undefined>>;
}
export declare const LinkTransferPreapproval:
  damlTypes.Template<LinkTransferPreapproval, undefined, '#link:Link.Token:LinkTransferPreapproval'> &
  damlTypes.ToInterface<LinkTransferPreapproval, never> &
  LinkTransferPreapprovalInterface

export declare type LockedLinkHolding = {
  lockedAdmin: damlTypes.Party,
  lockedSender: damlTypes.Party,
  lockedReceiver: damlTypes.Party,
  lockedInstrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  lockedAmount: damlTypes.Numeric,
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare interface LockedLinkHoldingInterface {
  Archive: 
    damlTypes.Choice<LockedLinkHolding, pkg9e70a8b3510d617f8a136213f33d6a903a10ca0eeec76bb06ba55d1ed9680f69.DA.Internal.Template.Archive, {}, undefined> &
    damlTypes.ChoiceFrom<damlTypes.Template<LockedLinkHolding, undefined>>;
}
export declare const LockedLinkHolding:
  damlTypes.Template<LockedLinkHolding, undefined, '#link:Link.Token:LockedLinkHolding'> &
  damlTypes.ToInterface<LockedLinkHolding, never> &
  LockedLinkHoldingInterface

export declare type Send = {
  sender: damlTypes.Party,
  amount: damlTypes.Numeric,
  instrumentId: pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.InstrumentId,
  inputHoldingCids: damlTypes.ContractId<pkg718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b.Splice.Api.Token.HoldingV1.Holding>[],
  meta: pkg4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f.Splice.Api.Token.MetadataV1.Metadata,
}

export declare const Send:
  damlTypes.Serializable<Send>
