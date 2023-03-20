package proto

type Msg struct {
	Ver       int    `json:"ver"`  // 版本
	Operation int    `json:"op"`   // 请求选项
	SeqId     string `json:"seq"`  // 客户端生成的序号
	Body      []byte `json:"body"` // 消息实体
}

type PushRoomMsgRequest struct {
	RoomId int
	Msg    Msg
}

type PushMsgRequest struct {
	UserId int
	Msg    Msg
}
