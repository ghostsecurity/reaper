import {RegExp} from "../regexp";
import {WorkflowM, NodeM, LinkM} from "../workflow";
import {HttpResponse, HttpRequest, KeyValue} from "../packaging";

export interface Workspace {
  id: string
  name: string
  scope: Scope
  interception_scope: Scope
  collection: Collection
  tree: Tree
  workflows: WorkflowM[]
}


export interface Collection {
  groups: Group[]
}


export interface Scope {
  include: Rule[]
  exclude: Rule[]
}


export interface Tree {
  root: StructureNode
}


export interface Rule {
  id: number
  protocol: string
  host: string
  path: string
  ports: number[]
}


export interface Request {
  id: string
  name: string
  inner: HttpRequest
  pre_script: string
  post_script: string
}


export interface StructureNode {
  id: string
  name: string
  children: StructureNode[]
}


export interface Group {
  id: string
  name: string
  requests: Request[]
}

