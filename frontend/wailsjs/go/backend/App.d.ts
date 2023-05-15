// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {workflow} from '../models';
import {workspace} from '../models';
import {settings} from '../models';
import {node} from '../models';

export function Confirm(arg1:string,arg2:string):Promise<boolean>;

export function CreateNode(arg1:number):Promise<workflow.NodeM>;

export function CreateWorkflow():Promise<workflow.WorkflowM>;

export function CreateWorkflowFromRequest(arg1:{[key: string]: any}):Promise<workflow.WorkflowM>;

export function CreateWorkspace(arg1:workspace.Workspace):Promise<workspace.Workspace>;

export function DeleteWorkspace(arg1:string):Promise<void>;

export function Error(arg1:string,arg2:string):Promise<void>;

export function ExportWorkflow(arg1:workflow.WorkflowM):Promise<void>;

export function GenerateID():Promise<string>;

export function GetSettings():Promise<settings.Settings>;

export function GetWorkspaces():Promise<Array<workspace.Workspace>>;

export function HighlightBody(arg1:string,arg2:string):Promise<string>;

export function HighlightHTTP(arg1:string):Promise<string>;

export function IgnoreThisUsedBindings(arg1:node.OutputM):Promise<workflow.UpdateM>;

export function ImportWorkflow():Promise<workflow.WorkflowM>;

export function LoadWorkspace(arg1:string):Promise<workspace.Workspace>;

export function Notify(arg1:string,arg2:string):Promise<void>;

export function RunWorkflow(arg1:workflow.WorkflowM):Promise<void>;

export function SaveSettings(arg1:settings.Settings):Promise<void>;

export function SaveWorkspace(arg1:workspace.Workspace):Promise<void>;

export function SelectFile(arg1:string):Promise<string>;

export function SetWorkspace(arg1:workspace.Workspace):Promise<void>;

export function StartProxy():Promise<void>;

export function StopProxy():Promise<void>;

export function StopWorkflow(arg1:workflow.WorkflowM):Promise<void>;

export function Warn(arg1:string,arg2:string):Promise<void>;
