import { VarStorageM, Connector } from '../node'

export interface NodeM {
  id: string
  name: string
  type: number
  vars: VarStorageM|null
  readonly: boolean
}

export interface Position {
  x: number
  y: number
}

export interface UpdateM {
  node: string
  status: string
  message: string
}

export interface WorkflowM {
  id: string
  name: string
  nodes: NodeM[]
  links: LinkM[]
  positioning: {[key: string]: Position}
}

export interface LinkM {
  from: LinkDirectionM
  to: LinkDirectionM
  annotation: string
}

export interface LinkDirectionM {
  node: string
  connector: string
}
