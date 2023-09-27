import {HttpRequest} from "../packaging";
import {WorkflowM, NodeM, LinkM} from "../workflow";

export interface StructureNode {
ID: string
Name: string
Children: StructureNode[]
}


export interface Group {
ID: string
Name: string
Requests: Request[]
}


export interface Request {
ID: string
Name: string
Inner: HttpRequest
PreScript: string
PostScript: string
}


export interface Tree {
Root: StructureNode
}


export interface Scope {
Include: Rule[]
Exclude: Rule[]
}


export interface Rule {
ID: number
Protocol: string
HostRegexRaw: string
HostRegex: Regexp|null
PathRegexRaw: string
PathRegex: Regexp|null
Ports: number[]
}


export interface Workspace {
ID: string
Name: string
Scope: Scope
InterceptionScope: Scope
Collection: Collection
Tree: Tree
Workflows: WorkflowM[]
}


export interface Collection {
Groups: Group[]
}

