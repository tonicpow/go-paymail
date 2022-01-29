package paymail

// Network an alias of the bitcoin networks.
type Network byte

const (
	// Mainnet bitcoin main network.
	Mainnet Network = iota
	// Testnet bitcoin test network.
	Testnet
	// STN bitcoin stress test network.
	STN
)

// String representation of the network.
func (n Network) String() string {
	switch n {
	case Mainnet:
		return "mainnet"
	case Testnet:
		return "testnet"
	case STN:
		return "STN"
	default:
		return "not recognized"
	}
}

// URLSuffix the conventional URL suffix for the network.
func (n Network) URLSuffix() string {
	switch n {
	case Testnet:
		return "-testnet"
	case STN:
		return "-stn"
	case Mainnet:
		return ""
	default:
		return ""
	}
}
