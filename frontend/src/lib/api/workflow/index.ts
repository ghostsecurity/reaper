import {VarStorageM} from "../node";

export interface NodeM {
Id: string
Name: string
Type: number
Vars: VarStorageM|null
ReadOnly: boolean
}


export interface LinkDirectionM {
Node: string
Connector: string
}


export interface Position {
X: number
Y: number
}


export interface WorkflowM {
ID: string
Name: string
Nodes: NodeM[]
Links: LinkM[]
Positioning: {[key: string]: Position}
}


export interface LinkM {
From: LinkDirectionM
To: LinkDirectionM
Annotation: string
}

