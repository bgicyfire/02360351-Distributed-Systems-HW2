// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.21.12
// source: multipaxos.proto

package multipaxos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TriggerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MemberId      string                 `protobuf:"bytes,1,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
	Event         *ScooterEvent          `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TriggerRequest) Reset() {
	*x = TriggerRequest{}
	mi := &file_multipaxos_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TriggerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerRequest) ProtoMessage() {}

func (x *TriggerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerRequest.ProtoReflect.Descriptor instead.
func (*TriggerRequest) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{0}
}

func (x *TriggerRequest) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

func (x *TriggerRequest) GetEvent() *ScooterEvent {
	if x != nil {
		return x.Event
	}
	return nil
}

type TriggerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ok            bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TriggerResponse) Reset() {
	*x = TriggerResponse{}
	mi := &file_multipaxos_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TriggerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerResponse) ProtoMessage() {}

func (x *TriggerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerResponse.ProtoReflect.Descriptor instead.
func (*TriggerResponse) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{1}
}

func (x *TriggerResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type PrepareRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Slot          int32                  `protobuf:"varint,1,opt,name=slot,proto3" json:"slot,omitempty"`
	Round         int32                  `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Id            string                 `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrepareRequest) Reset() {
	*x = PrepareRequest{}
	mi := &file_multipaxos_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareRequest) ProtoMessage() {}

func (x *PrepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareRequest.ProtoReflect.Descriptor instead.
func (*PrepareRequest) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{2}
}

func (x *PrepareRequest) GetSlot() int32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *PrepareRequest) GetRound() int32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *PrepareRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Promise message
type PrepareResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ok            bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	AcceptedRound int32                  `protobuf:"varint,2,opt,name=acceptedRound,proto3" json:"acceptedRound,omitempty"`      // Return previously accepted round
	AcceptedValue *ScooterEvent          `protobuf:"bytes,3,opt,name=acceptedValue,proto3,oneof" json:"acceptedValue,omitempty"` // Return previously accepted value
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrepareResponse) Reset() {
	*x = PrepareResponse{}
	mi := &file_multipaxos_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrepareResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareResponse) ProtoMessage() {}

func (x *PrepareResponse) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareResponse.ProtoReflect.Descriptor instead.
func (*PrepareResponse) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{3}
}

func (x *PrepareResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *PrepareResponse) GetAcceptedRound() int32 {
	if x != nil {
		return x.AcceptedRound
	}
	return 0
}

func (x *PrepareResponse) GetAcceptedValue() *ScooterEvent {
	if x != nil {
		return x.AcceptedValue
	}
	return nil
}

type AcceptRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Slot          int32                  `protobuf:"varint,1,opt,name=slot,proto3" json:"slot,omitempty"`
	Round         int32                  `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Id            string                 `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Value         *ScooterEvent          `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AcceptRequest) Reset() {
	*x = AcceptRequest{}
	mi := &file_multipaxos_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AcceptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptRequest) ProtoMessage() {}

func (x *AcceptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptRequest.ProtoReflect.Descriptor instead.
func (*AcceptRequest) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{4}
}

func (x *AcceptRequest) GetSlot() int32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *AcceptRequest) GetRound() int32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *AcceptRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AcceptRequest) GetValue() *ScooterEvent {
	if x != nil {
		return x.Value
	}
	return nil
}

// Learn message
type AcceptResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ok            bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AcceptResponse) Reset() {
	*x = AcceptResponse{}
	mi := &file_multipaxos_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AcceptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptResponse) ProtoMessage() {}

func (x *AcceptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptResponse.ProtoReflect.Descriptor instead.
func (*AcceptResponse) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{5}
}

func (x *AcceptResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type CommitRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Slot          int32                  `protobuf:"varint,1,opt,name=slot,proto3" json:"slot,omitempty"`
	Round         int32                  `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Id            string                 `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Value         *ScooterEvent          `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommitRequest) Reset() {
	*x = CommitRequest{}
	mi := &file_multipaxos_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitRequest) ProtoMessage() {}

func (x *CommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitRequest.ProtoReflect.Descriptor instead.
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{6}
}

func (x *CommitRequest) GetSlot() int32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *CommitRequest) GetRound() int32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *CommitRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CommitRequest) GetValue() *ScooterEvent {
	if x != nil {
		return x.Value
	}
	return nil
}

type CommitResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ok            bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommitResponse) Reset() {
	*x = CommitResponse{}
	mi := &file_multipaxos_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitResponse) ProtoMessage() {}

func (x *CommitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitResponse.ProtoReflect.Descriptor instead.
func (*CommitResponse) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{7}
}

func (x *CommitResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

// Define a generic Event message that can hold any type of event
type ScooterEvent struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	EventId   string                 `protobuf:"bytes,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	ScooterId string                 `protobuf:"bytes,2,opt,name=scooter_id,json=scooterId,proto3" json:"scooter_id,omitempty"` // Common field needed for all events
	// Types that are valid to be assigned to EventType:
	//
	//	*ScooterEvent_CreateEvent
	//	*ScooterEvent_ReserveEvent
	//	*ScooterEvent_ReleaseEvent
	EventType     isScooterEvent_EventType `protobuf_oneof:"event_type"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScooterEvent) Reset() {
	*x = ScooterEvent{}
	mi := &file_multipaxos_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScooterEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScooterEvent) ProtoMessage() {}

func (x *ScooterEvent) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScooterEvent.ProtoReflect.Descriptor instead.
func (*ScooterEvent) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{8}
}

func (x *ScooterEvent) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *ScooterEvent) GetScooterId() string {
	if x != nil {
		return x.ScooterId
	}
	return ""
}

func (x *ScooterEvent) GetEventType() isScooterEvent_EventType {
	if x != nil {
		return x.EventType
	}
	return nil
}

func (x *ScooterEvent) GetCreateEvent() *CreateScooterEvent {
	if x != nil {
		if x, ok := x.EventType.(*ScooterEvent_CreateEvent); ok {
			return x.CreateEvent
		}
	}
	return nil
}

func (x *ScooterEvent) GetReserveEvent() *ReserveScooterEvent {
	if x != nil {
		if x, ok := x.EventType.(*ScooterEvent_ReserveEvent); ok {
			return x.ReserveEvent
		}
	}
	return nil
}

func (x *ScooterEvent) GetReleaseEvent() *ReleaseScooterEvent {
	if x != nil {
		if x, ok := x.EventType.(*ScooterEvent_ReleaseEvent); ok {
			return x.ReleaseEvent
		}
	}
	return nil
}

type isScooterEvent_EventType interface {
	isScooterEvent_EventType()
}

type ScooterEvent_CreateEvent struct {
	CreateEvent *CreateScooterEvent `protobuf:"bytes,3,opt,name=create_event,json=createEvent,proto3,oneof"`
}

type ScooterEvent_ReserveEvent struct {
	ReserveEvent *ReserveScooterEvent `protobuf:"bytes,4,opt,name=reserve_event,json=reserveEvent,proto3,oneof"`
}

type ScooterEvent_ReleaseEvent struct {
	ReleaseEvent *ReleaseScooterEvent `protobuf:"bytes,5,opt,name=release_event,json=releaseEvent,proto3,oneof"`
}

func (*ScooterEvent_CreateEvent) isScooterEvent_EventType() {}

func (*ScooterEvent_ReserveEvent) isScooterEvent_EventType() {}

func (*ScooterEvent_ReleaseEvent) isScooterEvent_EventType() {}

// Specific event details for creating a scooter
type CreateScooterEvent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateScooterEvent) Reset() {
	*x = CreateScooterEvent{}
	mi := &file_multipaxos_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateScooterEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateScooterEvent) ProtoMessage() {}

func (x *CreateScooterEvent) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateScooterEvent.ProtoReflect.Descriptor instead.
func (*CreateScooterEvent) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{9}
}

// Specific event details for reserving a scooter
type ReserveScooterEvent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ReservationId string                 `protobuf:"bytes,1,opt,name=reservation_id,json=reservationId,proto3" json:"reservation_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReserveScooterEvent) Reset() {
	*x = ReserveScooterEvent{}
	mi := &file_multipaxos_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReserveScooterEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReserveScooterEvent) ProtoMessage() {}

func (x *ReserveScooterEvent) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReserveScooterEvent.ProtoReflect.Descriptor instead.
func (*ReserveScooterEvent) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{10}
}

func (x *ReserveScooterEvent) GetReservationId() string {
	if x != nil {
		return x.ReservationId
	}
	return ""
}

// Specific event details for releasing a scooter
type ReleaseScooterEvent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ReservationId string                 `protobuf:"bytes,1,opt,name=reservation_id,json=reservationId,proto3" json:"reservation_id,omitempty"`
	Distance      int64                  `protobuf:"varint,2,opt,name=distance,proto3" json:"distance,omitempty"` // Distance the scooter was used for
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReleaseScooterEvent) Reset() {
	*x = ReleaseScooterEvent{}
	mi := &file_multipaxos_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleaseScooterEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleaseScooterEvent) ProtoMessage() {}

func (x *ReleaseScooterEvent) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleaseScooterEvent.ProtoReflect.Descriptor instead.
func (*ReleaseScooterEvent) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{11}
}

func (x *ReleaseScooterEvent) GetReservationId() string {
	if x != nil {
		return x.ReservationId
	}
	return ""
}

func (x *ReleaseScooterEvent) GetDistance() int64 {
	if x != nil {
		return x.Distance
	}
	return 0
}

var File_multipaxos_proto protoreflect.FileDescriptor

var file_multipaxos_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x22, 0x5a, 0x0a, 0x0e, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x63, 0x6f, 0x6f,
	0x74, 0x65, 0x72, 0x2e, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x21, 0x0a, 0x0f, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x4a, 0x0a, 0x0e, 0x50, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x9b, 0x01, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x70, 0x61,
	0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64,
	0x12, 0x40, 0x0a, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65,
	0x72, 0x2e, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00,
	0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88,
	0x01, 0x01, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x76, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x2b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x20, 0x0a, 0x0e,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x76,
	0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73,
	0x6c, 0x6f, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74,
	0x65, 0x72, 0x2e, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x20, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0xa2, 0x02, 0x0a, 0x0c, 0x53, 0x63, 0x6f,
	0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x40, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x73, 0x63, 0x6f, 0x6f,
	0x74, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65,
	0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x43, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73,
	0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x53, 0x63,
	0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0c, 0x72, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x43, 0x0a, 0x0d, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6c, 0x65,
	0x61, 0x73, 0x65, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48,
	0x00, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42,
	0x0c, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x14, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x22, 0x3c, 0x0a, 0x13, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x53, 0x63,
	0x6f, 0x6f, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x22, 0x58, 0x0a, 0x13, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x53, 0x63, 0x6f, 0x6f,
	0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x65,
	0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x32, 0x91, 0x02, 0x0a, 0x11,
	0x4d, 0x75, 0x6c, 0x74, 0x69, 0x50, 0x61, 0x78, 0x6f, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x42, 0x0a, 0x0d, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x4c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x12, 0x17, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x63,
	0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x12, 0x17, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61,
	0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x63, 0x6f, 0x6f,
	0x74, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x06, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12,
	0x16, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65,
	0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3b, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x2e, 0x73,
	0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x67,
	0x69, 0x63, 0x79, 0x66, 0x69, 0x72, 0x65, 0x2f, 0x30, 0x32, 0x33, 0x36, 0x30, 0x33, 0x35, 0x31,
	0x2d, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x2d, 0x53, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x73, 0x2d, 0x48, 0x57, 0x32, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_multipaxos_proto_rawDescOnce sync.Once
	file_multipaxos_proto_rawDescData = file_multipaxos_proto_rawDesc
)

func file_multipaxos_proto_rawDescGZIP() []byte {
	file_multipaxos_proto_rawDescOnce.Do(func() {
		file_multipaxos_proto_rawDescData = protoimpl.X.CompressGZIP(file_multipaxos_proto_rawDescData)
	})
	return file_multipaxos_proto_rawDescData
}

var file_multipaxos_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_multipaxos_proto_goTypes = []any{
	(*TriggerRequest)(nil),      // 0: scooter.TriggerRequest
	(*TriggerResponse)(nil),     // 1: scooter.TriggerResponse
	(*PrepareRequest)(nil),      // 2: scooter.PrepareRequest
	(*PrepareResponse)(nil),     // 3: scooter.PrepareResponse
	(*AcceptRequest)(nil),       // 4: scooter.AcceptRequest
	(*AcceptResponse)(nil),      // 5: scooter.AcceptResponse
	(*CommitRequest)(nil),       // 6: scooter.CommitRequest
	(*CommitResponse)(nil),      // 7: scooter.CommitResponse
	(*ScooterEvent)(nil),        // 8: scooter.ScooterEvent
	(*CreateScooterEvent)(nil),  // 9: scooter.CreateScooterEvent
	(*ReserveScooterEvent)(nil), // 10: scooter.ReserveScooterEvent
	(*ReleaseScooterEvent)(nil), // 11: scooter.ReleaseScooterEvent
}
var file_multipaxos_proto_depIdxs = []int32{
	8,  // 0: scooter.TriggerRequest.event:type_name -> scooter.ScooterEvent
	8,  // 1: scooter.PrepareResponse.acceptedValue:type_name -> scooter.ScooterEvent
	8,  // 2: scooter.AcceptRequest.value:type_name -> scooter.ScooterEvent
	8,  // 3: scooter.CommitRequest.value:type_name -> scooter.ScooterEvent
	9,  // 4: scooter.ScooterEvent.create_event:type_name -> scooter.CreateScooterEvent
	10, // 5: scooter.ScooterEvent.reserve_event:type_name -> scooter.ReserveScooterEvent
	11, // 6: scooter.ScooterEvent.release_event:type_name -> scooter.ReleaseScooterEvent
	0,  // 7: scooter.MultiPaxosService.TriggerLeader:input_type -> scooter.TriggerRequest
	2,  // 8: scooter.MultiPaxosService.Prepare:input_type -> scooter.PrepareRequest
	4,  // 9: scooter.MultiPaxosService.Accept:input_type -> scooter.AcceptRequest
	6,  // 10: scooter.MultiPaxosService.Commit:input_type -> scooter.CommitRequest
	1,  // 11: scooter.MultiPaxosService.TriggerLeader:output_type -> scooter.TriggerResponse
	3,  // 12: scooter.MultiPaxosService.Prepare:output_type -> scooter.PrepareResponse
	5,  // 13: scooter.MultiPaxosService.Accept:output_type -> scooter.AcceptResponse
	7,  // 14: scooter.MultiPaxosService.Commit:output_type -> scooter.CommitResponse
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_multipaxos_proto_init() }
func file_multipaxos_proto_init() {
	if File_multipaxos_proto != nil {
		return
	}
	file_multipaxos_proto_msgTypes[3].OneofWrappers = []any{}
	file_multipaxos_proto_msgTypes[8].OneofWrappers = []any{
		(*ScooterEvent_CreateEvent)(nil),
		(*ScooterEvent_ReserveEvent)(nil),
		(*ScooterEvent_ReleaseEvent)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_multipaxos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_multipaxos_proto_goTypes,
		DependencyIndexes: file_multipaxos_proto_depIdxs,
		MessageInfos:      file_multipaxos_proto_msgTypes,
	}.Build()
	File_multipaxos_proto = out.File
	file_multipaxos_proto_rawDesc = nil
	file_multipaxos_proto_goTypes = nil
	file_multipaxos_proto_depIdxs = nil
}
