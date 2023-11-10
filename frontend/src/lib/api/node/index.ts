export interface VarStorageM {
  inputs: Connector[]
  outputs: Connector[]
  static: {[key: string]: TransmissionM}
}

export interface OutputM {
  node: string
  channel: string
  message: string
}

export interface TransmissionM {
  type: number
  internal: number
  data: any
}

export interface Connector {
  name: string
  type: number
  linkable: boolean
  description: string
}
