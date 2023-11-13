import { VarStorageM } from '../node'

export interface UpdateM {
  node: string;
  status: string;
  message: string;
}

export interface NodeM {
  id: string;
  name: string;
  type: number;
  vars: VarStorageM|null;
  readonly: boolean;
}

export interface LinkDirectionM {
  node: string;
  connector: string;
}

export interface LinkM {
  from: LinkDirectionM;
  to: LinkDirectionM;
  annotation: string;
}

export interface Position {
  x: number;
  y: number;
}

export interface WorkflowM {
  id: string;
  version: number;
  name: string;
  nodes: NodeM[];
  links: LinkM[];
  positioning: {[key: string]: Position};
}

