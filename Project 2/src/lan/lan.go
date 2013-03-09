package lan 

import (
	"stats"
)

// Holds a Start and End "time", or Tick.
// These are when the packet starts arriving at a node, and when it finishes.
// A 0 value indicates that there is no packet set to arrive at the node.
type Packet_Arrival struct {
	Start	float64
	End   	float64
	p 		*stats.Packet  
}

var (

	TICK_time 	float64 // 1 TICK = X milliseconds

	num_comps 		 	float64 // Number of computers  connected tot he LAN
	packets_per_sec 	float64 // Number of packets/sec
	W		  			float64 // The speed of the lan in bits/msec
	L         			float64 // Length of a packet in bits

	sense_flags 	[]bool				// Flags of whether or not data is arrive at that node during the current tick.
										// Each element correspnds to a computer. Ie: comp0  is at element 0, comp1 and element 1.

	node_info []Packet_Arrival			// The Packet_Arrival struct for each computer.
										// Each element correspnds to a computer.

	Prop_Ticks				float64
	Packer_Trans_Ticks		float64
	Jam_Trans_Ticks			float64
)	



// Creates the two slices: sense_flags, node_info each with length equal to the number of computers
// Initializes all relevant input variables.

func Init(tick_t float64, n float64, a float64, w float64, l float64) {
	TICK_time = tick_t;
	num_comps = n;
	packets_per_sec = a;
	W = w;
	L = l;

	sense_flags = make([]bool, N, N)
	node_info = make([]Packet_Arrival, N, N)

	// will be calculated, not hard coded.
	Prop_Ticks = 10
	Packet_Trans_Ticks = 100
	Jam_Trans_Ticks = 50

}

// For each "computer" check if the current tick is betweeen the start and end times of 
// when a packet is arriving at that node.
// If it is between them, sense_flags for that node becomes true.
// If t = the End time, then the packet has just fully arrived and sense_flags = false. THe packet_arrival struct is also cleared.
func Complete_Tick(t float64) {
	for i:= 0; i < num_comps; i++ {
		// Packet is currently arriving at this node
		if t >= node_info[i].Start && t < node_info[i].End {
			sense_flags[i] = true
		}
		// packet has just fully arrived at this node.
		// clear Packet_Arrival struct.
		if t == node_info[i].End {
			sense_flags[i] = false
			node_info[i].Start = 0
			node_info[i].End = 0
			push_to_bucket(node_info[i].p)
			node_info[i].p = nil
		}

	}
}

// pushes the packet to the bucket.
func push_to_bucket(p *stats.Packet)
{

}

// returns whether there is currently data arriving at the specified computer.
func sense_line(compID float64) bool {
	return sense_flags[compID];
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
func put_packet(p *Packet, compID float64, Current_Tick float64) bool {

	Start := Current_Tick + Prop_Ticks
	End := Start + Packet_Trans_Ticks

	for i := 0; i < num_comps; i++ {
		if i != compID {
			if End > node_info[i].End {
				node_info[i].End = End
				node_info[i].p = p
			}
			if node_info[i].Start == 0 {
				node_info[i].Start = Start
			}
		}
	}
}

// Has similar functionality to put_packet, except there is no actual packet to keep track of, only a jam signal.
// The difference is that once a jam signal is sent, The Packet_Arrival.End now becomes:
// 		current tick + Jam_Trans_Ticks + Jam_Prop_Ticks;
// This is because if a computer is sending a jam signal, then once it arrives at other nodes, they stop transmitting packets.
// At this point, only Jam signals will be on the line. Once the last Jam signal has fully arrived, the line has to be open. 
func send_jam_signal(compID float64, Current_Tick float64) bool {

	Start := Current_Tick + Prop_Ticks;
	End := Start + Jam_Trans_Ticks;

	for i := 0; i < num_comps; i++ {
		if i != compID {
			node_info[i].End = End
			node_info[i].p = nil
			if node_info[i].Start == 0 {
			    node_info[i].Start = Start
			}
		}
	}
}





