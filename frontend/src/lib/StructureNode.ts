
export default class StructureNode {
    Name: string;
    Children: Array<StructureNode>;
    constructor(name: string, children: StructureNode[]) {
        this.Name = name;
        this.Children = children;
    }
}