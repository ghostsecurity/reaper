import { HttpRequest } from '../packaging'
import { WorkflowM } from '../workflow'

export interface Rule {
  id: number;
  protocol: string;
  host: string;
  path: string;
  ports: number[];
}

export interface Scope {
  include: Rule[];
  exclude: Rule[];
}

export interface Request {
  id: string;
  name: string;
  inner: HttpRequest;
  pre_script: string;
  post_script: string;
}

export interface Group {
  id: string;
  name: string;
  requests: Request[];
}

export interface Collection {
  groups: Group[];
}

export interface StructureNode {
  id: string;
  name: string;
  children: StructureNode[];
}

export interface Tree {
  root: StructureNode;
}

export interface Workspace {
  id: string;
  name: string;
  scope: Scope;
  interception_scope: Scope;
  collection: Collection;
  tree: Tree;
  workflows: WorkflowM[];
}
