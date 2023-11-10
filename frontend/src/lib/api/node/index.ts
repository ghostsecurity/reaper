import { ParentType } from '../transmission'

export interface OutputM {
  node: string;
  channel: string;
  message: string;
}

export interface Connector {
  name: string;
  type: ParentType;
  linkable: boolean;
  description: string;
}

export interface TransmissionM {
  type: number;
  internal: number;
  data: any;
}

export interface VarStorageM {
  inputs: Connector[];
  outputs: Connector[];
  static: {[key: string]: TransmissionM};
}
