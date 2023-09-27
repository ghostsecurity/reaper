
export interface VarStorageM {
Inputs: Connector[]
Outputs: Connector[]
Static: {[key: string]: TransmissionM}
}


export interface TransmissionM {
ParentType: number
ChildType: number
Data: any
}


export interface Connector {
Name: string
Type: ParentType
Linkable: boolean
Description: string
}

