/*
    WARNING:
        This file is automatically generated. Do not edit it directly!
*/
import { MessageType, WebsocketMessage } from './websocket.type'
// %IMPORTS:START%
import { UpdateM, NodeM, WorkflowM } from './workflow'
import { OutputM } from './node'
import { Workspace } from './workspace'
import { HttpRequest } from './packaging'
import { Settings } from './settings'
import { VersionInfo } from './api'
// %IMPORTS:END%

export default class Client {
  ws: WebSocket | null = null

  awaitingRes: Map<string, (args: string[]) => void>

  awaitingRej: Map<string, (reason?: string) => void>

  eventCallbacks: Map<string, ((...args: any) => void)[]>

  index: number

  url: string

  constructor(url?: string) {
        this.awaitingRes = new Map<string, (args: string[]) => void>() //eslint-disable-line
        this.awaitingRej = new Map<string, (reason?: string) => void>() //eslint-disable-line
        this.eventCallbacks = new Map<string, ((o1?: any, o2?: any, o3?: any) => void)[]>() //eslint-disable-line
    this.index = 0
    if (url) {
      this.url = url
    } else if (typeof window !== 'undefined') {
      this.url = `ws://${window.location.host}/ws/`
    } else {
      this.url = 'ws://127.0.0.1:31337/ws/'
    }
  }

  Init(): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      let ready = false
      this.ws = new WebSocket(this.url)
      this.ws.onerror = err => {
        if (!ready) {
          reject(err)
        }
      }
      this.ws.onclose = () => {
        this.eventCallbacks.get('Close')?.forEach(callback => {
          callback.apply(null, [])
        })
      }
      this.ws.onopen = () => this.sendMessage({
        messageType: MessageType.Ping,
        identifier: 'client',
        args: [],
        sender: '',
      })
      this.ws.onmessage = data => {
        const message = <WebsocketMessage>JSON.parse(data.data)
        switch (message.messageType) {
          case MessageType.Ping:
            this.sendMessage({
              messageType: MessageType.Pong,
              identifier: message.identifier,
              args: [],
              sender: '',
            })
            break
          case MessageType.Pong:
            if (!ready) {
              ready = true
              resolve(true)
            }
            break
          case MessageType.Result:
            // eslint-disable-next-line no-case-declarations
            const prevResolve = this.awaitingRes.get(message.sender)
            if (prevResolve != null) {
              prevResolve.call(null, message.args)
              this.awaitingRes.delete(message.sender)
              this.awaitingRej.delete(message.sender)
            }
            break
          case MessageType.Failure:
            // eslint-disable-next-line no-case-declarations
            const prevReject = this.awaitingRej.get(message.sender)
            if (prevReject != null) {
              const err: string = JSON.parse(message.args[0])
              prevReject.call(null, err)
              this.awaitingRes.delete(message.sender)
              this.awaitingRej.delete(message.sender)
            }
            break
          case MessageType.Notify:
            // fire callback(s)
            this.eventCallbacks.get(message.identifier)?.forEach(callback => {
              const actualArgs = []
              for (let i = 0; i < message.args.length; i += 1) {
                actualArgs.push(JSON.parse(message.args[i]))
              }
              callback.apply(null, actualArgs)
            })
            break
          case MessageType.Error:
            if (!ready) {
              const err: string = JSON.parse(message.args[0])
              reject(err)
            }
            break
          default:
            console.log(message)
            break
        }
      }
    })
  }

  close() {
    if (this.ws == null) {
      return
    }
    this.ws.close()
    this.ws = null
  }

  sendMessage(message: WebsocketMessage) {
    if (this.ws == null) {
      return
    }
    this.ws.send(JSON.stringify(message))
  }

  OnEvent(eventName: string, callback: (o1?: any, o2?: any, o3?: any) => void): void {
    let current = this.eventCallbacks.get(eventName)
    if (!current) {
      current = []
    }
    current.push(callback)
    this.eventCallbacks.set(eventName, current)
    this.sendMessage({
      messageType: MessageType.Subscribe,
      identifier: eventName,
    } as WebsocketMessage)
  }

  callMethod(method: string, args: any[], receive: (args: string[]) => void, reject: (reason?: any) => void) {
    const sender = this.index.toString()
    this.index += 1
    this.awaitingRes.set(sender, receive)
    this.awaitingRej.set(sender, reject)
    for (let i = 0; i < args.length; i += 1) {
      args[i] = JSON.stringify(args[i])
    }
    this.sendMessage({
      messageType: MessageType.Method,
      identifier: method,
      args,
      sender,
    })
  }

  // %METHODS:START%
  BindingsOnly(a0: UpdateM, a1: OutputM): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('BindingsOnly', [a0, a1], receive, reject)
    })
  }

  Close(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('Close', [], receive, reject)
    })
  }

  CreateNode(a0: number): Promise<NodeM|null> {
    return new Promise<NodeM|null>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: NodeM|null = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('CreateNode', [a0], receive, reject)
    })
  }

  CreateWorkflow(): Promise<WorkflowM|null> {
    return new Promise<WorkflowM|null>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: WorkflowM|null = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('CreateWorkflow', [], receive, reject)
    })
  }

  CreateWorkflowFromRequest(a0: {[key: string]: any}): Promise<WorkflowM|null> {
    return new Promise<WorkflowM|null>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: WorkflowM|null = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('CreateWorkflowFromRequest', [a0], receive, reject)
    })
  }

  CreateWorkspace(a0: Workspace|null): Promise<Workspace|null> {
    return new Promise<Workspace|null>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: Workspace|null = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('CreateWorkspace', [a0], receive, reject)
    })
  }

  DeleteWorkspace(a0: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('DeleteWorkspace', [a0], receive, reject)
    })
  }

  DropInterceptedRequest(a0: HttpRequest): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('DropInterceptedRequest', [a0], receive, reject)
    })
  }

  FormatCode(a0: string, a1: string): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: string = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('FormatCode', [a0, a1], receive, reject)
    })
  }

  GenerateID(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: string = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('GenerateID', [], receive, reject)
    })
  }

  GetSettings(): Promise<Settings> {
    return new Promise<Settings>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: Settings = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('GetSettings', [], receive, reject)
    })
  }

  GetVersionInfo(): Promise<VersionInfo> {
    return new Promise<VersionInfo>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: VersionInfo = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('GetVersionInfo', [], receive, reject)
    })
  }

  GetWorkspace(): Promise<Workspace|null> {
    return new Promise<Workspace|null>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: Workspace|null = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('GetWorkspace', [], receive, reject)
    })
  }

  GetWorkspaces(): Promise<Workspace[]> {
    return new Promise<Workspace[]>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: Workspace[] = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('GetWorkspaces', [], receive, reject)
    })
  }

  HighlightBody(a0: string, a1: string): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: string = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('HighlightBody', [a0, a1], receive, reject)
    })
  }

  HighlightHTTP(a0: string): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: string = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('HighlightHTTP', [a0], receive, reject)
    })
  }

  LoadWorkspace(a0: string): Promise<Workspace|null> {
    return new Promise<Workspace|null>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: Workspace|null = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('LoadWorkspace', [a0], receive, reject)
    })
  }

  ModifyInterceptedRequest(a0: HttpRequest): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('ModifyInterceptedRequest', [a0], receive, reject)
    })
  }

  RunWorkflow(a0: WorkflowM|null): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('RunWorkflow', [a0], receive, reject)
    })
  }

  SaveSettings(a0: Settings|null): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('SaveSettings', [a0], receive, reject)
    })
  }

  SaveWorkspace(a0: Workspace|null): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('SaveWorkspace', [a0], receive, reject)
    })
  }

  SendRequest(a0: HttpRequest): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('SendRequest', [a0], receive, reject)
    })
  }

  SetWorkspace(a0: Workspace|null): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('SetWorkspace', [a0], receive, reject)
    })
  }

  StartProxy(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('StartProxy', [], receive, reject)
    })
  }

  StopProxy(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('StopProxy', [], receive, reject)
    })
  }

  StopWorkflow(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const receive = () => {
        resolve()
      }
      this.callMethod('StopWorkflow', [], receive, reject)
    })
  }

  Test(a0: string): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      const receive = (args: string[]) => {
        const output0: string = JSON.parse(args[0])
        resolve(output0)
      }
      this.callMethod('Test', [a0], receive, reject)
    })
  }

  // %METHODS:END%

  /*
                          Test(a: string): Promise<string> {
                          let res: (value: string | PromiseLike<string>) => void;
                          let rej: (reason?: any) => void;
                          return new Promise<string>((resolve, reject) => {
                              res = resolve;
                              rej = reject;
                              const receive = (_: string[]) => {
                                  let a: string = JSON.parse(args[0]);
                                  res(a);
                              }
                              this.callMethod("Test", [a], receive, rej);
                          })

                       */
}
