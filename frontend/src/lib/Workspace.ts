
interface Workspace {
    ID: string;
    Name: string;
    Scope : Scope;
}

interface Scope {
    Include: Array<Rule>;
    Exclude: Array<Rule>;
}

interface Rule {
    Protocol:  string;
    HostRegex: string;
    PathRegex: string;
    Ports:     Array<number>;
}

export type {
    Workspace,
    Scope,
    Rule
}
