package core

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type Dispatch struct {
	device string
	Plug   *Plug
}

func NewDispatch(plug *Plug, device string) *Dispatch {
	return &Dispatch{
		Plug:   plug,
		device: device,
	}
}

func (d *Dispatch) Capture() {
	// init device
	handle, err := pcap.OpenLive(d.device, 65535, false, pcap.BlockForever)
	if err != nil {
		fmt.Println(err)
		return
	}

	// set filter
	fmt.Println(d.Plug.BPF)
	err = handle.SetBPFFilter(d.Plug.BPF)
	if err != nil {
		fmt.Println(err)
	}

	// capture
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := src.Packets()

	// set up assembly
	streamFactory := &ProtocolStreamFactory{
		dispatch: d,
	}
	streamPool := NewStreamPool(streamFactory)
	assembler := NewAssembler(streamPool)
	ticker := time.NewTicker(time.Minute)

	// loop until ctrl+z
	for {
		select {
		case packet := <-packets:
			if packet.NetworkLayer() == nil ||
				packet.TransportLayer() == nil ||
				packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				fmt.Println("ERR : Unknown Packet -_-")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(
				packet.NetworkLayer().NetworkFlow(),
				tcp, packet.Metadata().Timestamp,
			)
		case <-ticker.C:
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}

type ProtocolStreamFactory struct {
	dispatch *Dispatch
}

type ProtocolStream struct {
	r tcpreader.ReaderStream
}

func (m *ProtocolStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	// init stream struct
	stm := &ProtocolStream{
		r: tcpreader.NewReaderStream(),
	}

	// new stream
	fmt.Println("# Start new stream:", net, transport)

	// decode packet
	go m.dispatch.Plug.ResolveStream(net, transport, &(stm.r))

	return &(stm.r)
}
