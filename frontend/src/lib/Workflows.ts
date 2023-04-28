export const enum NodeType {
    UNKNOWN = 0,
    FUZZER = 1,
    OUTPUT = 2,
    STATUS_FILTER = 3,
    REQUEST = 4,
    START = 5,
    SENDER = 6,
}

export const enum ParentType {
    UNKNOWN = 1,
    STRING = 2,
    INT = 4,
    MAP = 8,
    LIST = 16,
    REQUEST = 32,
    RESPONSE = 64,
    START = 128,
    BOOLEAN = 256,
}

export const enum ChildType {
    UNKNOWN = 1,
    NUMERIC_RANGE_LIST = 2,
    WORD_LIST = 3,
}

export function NodeTypeName(t: NodeType): string {
    switch (t) {
        case NodeType.FUZZER:
            return 'Fuzzer'
        case NodeType.STATUS_FILTER:
            return 'Status Filter'
        case NodeType.OUTPUT:
            return 'Output'
        case NodeType.REQUEST:
            return 'Request'
        case NodeType.START:
            return 'Start'
        case NodeType.SENDER:
            return 'Sender'
        default:
            return 'Unknown (' + t + ')'
    }
}
