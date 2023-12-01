package main

import (
	"log"
	"pb/pb"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//	func TestMessageIncorrectType(t *testing.T) {
//		println(t.Name())
//		order := exampleMessage()
//		// wraps Order in ProtoBufs higher order Message type for easy conversion back to Order
//		any, _ := anypb.New(order)
//		byteSlice, _ := proto.Marshal(any)
//		// turn byteSlice back to any
//		newAny := &anypb.Any{}
//		proto.Unmarshal(byteSlice, newAny)
//		// create new Order message type to be converted from Any
//		incorrectType := &pb.OrderLine{}
//		// replace unMarshaledOrder to what any is
//		err := newAny.UnmarshalTo(incorrectType)
//		if err != nil {
//			println(err)
//		}
//	}
//
//	func TestMessageTypeCheck(t *testing.T) {
//		println(t.Name())
//		order := exampleMessage()
//		// wraps Order in ProtoBufs higher order Message type for easy conversion back to Order
//		any, _ := anypb.New(order)
//		byteSlice, _ := proto.Marshal(any)
//		// turn byteSlice back to any
//		newAny := &anypb.Any{}
//		proto.Unmarshal(byteSlice, newAny)
//		// create new Order message type to be converted from Any
//		newOrder, err := newAny.UnmarshalNew()
//		if err != nil {
//			println(err)
//		} else {
//			println(protojson.Format(newOrder))
//		}
//	}
func TestPushMessage(t *testing.T) {
	println(t.Name())

	// setup server
	cfg := &Config{
		ListenAddr:        ":3000",
		QueueProducerFunc: QueueFunc("Name"),
	}
	s, _ := NewServer(cfg)
	go s.Start()

	// creat the topic
	topic := "Order"
	err := s.CreateTopic(topic)
	if err != nil {
		log.Fatal(err)
	}

	// create the message
	order := exampleMessage()
	any, _ := anypb.New(order)
	byteSlice, _ := proto.Marshal(any)

	// push the message and get the offset
	offset, err := s.topics[topic].Push(byteSlice)
	if err != nil {
		log.Fatal(err)
	}

	// retrieve the message with the offset
	retrievedByteSlice, err := s.topics[topic].Grab(offset)
	if err != nil {
		log.Fatal(err)
	}

	// un serialse the message
	newAny := &anypb.Any{}
	proto.Unmarshal(retrievedByteSlice, newAny)
	newOrder, err := newAny.UnmarshalNew()

	if err != nil {
		println(err)
	} else {
		println(protojson.Format(newOrder))
	}
}

func exampleMessage() *pb.CustomerOrder {
	line := &pb.OrderLine{
		ProductNumber: 12345,
		ProductName:   "Socks",
		Qty:           2,
	}
	lines := make([]*pb.OrderLine, 0)
	lines = append(lines, line)
	order := &pb.CustomerOrder{
		CustomerNumber:  1,
		CustomerName:    "Winston",
		CustomerAddress: "123 Fake Street",
		OrderTime:       timestamppb.Now(),
		ShipDate:        timestamppb.New(time.Date(2023, 12, 25, 12, 0, 0, 0, time.Local)),
		Lines:           lines,
	}
	return order
}
