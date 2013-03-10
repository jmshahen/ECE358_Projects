package lan 

import (
	"stats"
)

// Holds a Start and End "time", or Tick.
// These are when the packet starts arriving at a node, and when it finishes.
// A 0 value indicates that there is no packet set to arrive at the node.
type Packet_Arrival struct {
	Start	int64
	End   	int64
	p 		*stats.Packet  
}

type LAN struct {

	current_tick 			int64

	num_comps 		 		int64 	// Number of computers  connected tot he LAN	
	Prop_Ticks				int64
	Packer_Trans_Ticks		int64
	Jam_Trans_Ticks			int64

	sense_flags 			[]bool		// Flags of whether or not data is arrive at that node during the current tick.
										// Each element correspnds to a computer. Ie: comp0  is at element 0, comp1 and element 1.

	node_info []Packet_Arrival			// The Packet_Arrival struct for each computer.
										// Each element correspnds to a computer.
	bucket 					*stats.Bucket
	lost_bucket				*stats.Bucket
}


// Creates the two slices: sense_flags, node_info each with length equal to the number of computers
// Initializes all relevant input variables.


func (LAN *lan) Init(n int64, prop_t int64, pack_t int64, jam_t int64 64, b *stats.Bucket, l_b *stats.Bucket) {

	lan.sense_flags = make([]bool, N, N)
	lan.node_info = make([]Packet_Arrival, N, N)

	lan.num_comps = n;
	lan.Prop_Ticks = prop_t
	lan.Packet_Trans_Ticks = pack_t
	lan.Jam_Trans_Ticks = jam_t;

	lan.bucket = b
	lan.lost_bucket = l_b

}



// For each "computer" check if the current tick is betweeen the start and end times of 
// when a packet is arriving at that node.
// If it is between them, sense_flags for that node becomes true.
// If t = the End time, then the packet has just fully arrived and sense_flags = false. THe packet_arrival struct is also cleared.
func (LAN *lan) Complete_Tick(t int64) {
	lan.current_tick = t;
	for i:= 0; i < lan.num_comps; i++ {
		// Packet is currently arriving at this node
		if t >= lan.node_info[i].Start && t < lan.node_info[i].End {
			lan.sense_flags[i] = true
		}
		// packet has just fully arrived at this node.
		// clear Packet_Arrival struct.
		if t == lan.node_info[i].End {
			lan.sense_flags[i] = false
			lan.node_info[i].Start = 0
			lan.node_info[i].End = 0
			lan.push_to_bucket(node_info[i].p)
			lan.node_info[i].p = nil
		}
	}
}

// pushes the packet to the bucket.
func (LAN *lan) push_to_bucket(p *stats.Packet)
{
		p.Finished = lan.current_tick
		lan.bucket.Accept_packet(p)
}

func (LAN *lan) record_lost_packet(p *stats.Packet)
{
		p.Finished = lan.current_tick
		lan.lost_bucket.Accept_packet(packet)
}

// returns whether there is currently data arriving at the specified computer.
func (LAN *lan) sense_line(compID int64) bool {
	return lan.sense_flags[compID];
}

// Called by a computer to put a packet on the line.
// Should only be called once the computer has been sensing the line for the appropriate amount of time.
// Once the packet is put on the line, the node_info slice is updated.
// For each computer node except the sending node, the Start and End variables of the Packet_Arrival struct for that node
// are updated with the correct tick times.
//
// Ex:
// Computer 3 puts a packet on the line at tick 500.
// Transmission time for a packet is 100 ticks
// the progagation time is 5 ticks for it to get to all nodes.
// The packet then arrives at all the other nodes starting at tick 505, and has fully arrived at tick 605.
// For each Computer node except node 3, the Start field of Packet_Arrival is set to 505, and the End field to 605.
//
// If for example, Node 7 starts transmitting at Tick 502, with all other properties similar then:
// Its packet_arrival struct will remain unchanged
// For Node 3, Start = 507, End = 607
// For the rest of the nodes, the first packet will have fully arrived at 505, but the secound packet will have
// not fully arrived until 607. So: Start = 505, End = 607. 
//
// We always use the greater of Packet_Arrival.End and our calculated End.
// We only change the Packet_Arrival.Start field if it is 0.
// If it is not zero then this means that it has to be less than our calculated start, and we should not replace it.
func (LAN *lan) put_packet(p *Packet, compID int64, Current_Tick int64) int64 {
	Start := Current_Tick + lan.Prop_Ticks
	End := Start + lan.Packet_Trans_Ticks

	for i := 0; i < lan.num_comps; i++ {
		if i != compID {
			if End > lan.node_info[i].End {
				lan.node_info[i].End = End
				lan.node_info[i].p = p
			}
			if lan.node_info[i].Start == 0 {
				lan.node_info[i].Start = Start
			}
		}
	}
	return lan.Packet_Trans_Ticks
}

// Has similar functionality to put_packet, except there is no actual packet to keep track of, only a jam signal.
// The difference is that once a jam signal is sent, The Packet_Arrival.End now becomes:
// 		current tick + Jam_Trans_Ticks + Jam_Prop_Ticks;
// This is because if a computer is sending a jam signal, then once it arrives at other nodes, they stop transmitting packets.
// At this point, only Jam signals will be on the line. Once the last Jam signal has fully arrived, the line has to be open. 
func (LAN *lan) send_jam_signal(compID int64, Current_Tick int64) int64 {
	Start := Current_Tick + lan.Prop_Ticks;
	End := Start + lan.Jam_Trans_Ticks;

	for i := 0; i < lan.num_comps; i++ {
		if i != compID {
			lan.node_info[i].End = End
			lan.node_info[i].p = nil
			if lan.node_info[i].Start == 0 {
			    lan.node_info[i].Start = Start
			}
		}
	}
	return lan.Packet_Prop_Ticks
}





