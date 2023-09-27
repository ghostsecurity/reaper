export type WebsocketMessage = {
    messageType: MessageType
    identifier: string
    args: string[]
    sender: string
};

export const enum MessageType {
    Unknown = 0,
    Ping = 1,
    Pong = 2,
    Subscribe = 3,
    Notify = 4,
    Method = 5,
    Result = 6,
    Failure = 7,
    Error = 8
}
