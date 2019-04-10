package nine

import (
	"strings"
	"time"
)

type Mapstringstring map[string]*string

func (m Mapstringstring) String() (out string) {
	for i, x := range m {
		out += i + ":" + *x + " "
	}
	return strings.TrimSpace(out)
}

type Config struct {
	ConfigFile               *string
	AppDataDir               *string
	DataDir                  *string
	LogDir                   *string
	LogLevel                 *string
	Subsystems               *Mapstringstring
	Network                  *string
	AddPeers                 *[]string
	ConnectPeers             *[]string
	MaxPeers                 *int
	Listeners                *[]string
	DisableListen            *bool
	DisableBanning           *bool
	BanDuration              *time.Duration
	BanThreshold             *int
	Whitelists               *[]string
	Username                 *string
	Password                 *string
	ServerUser               *string
	ServerPass               *string
	LimitUser                *string
	LimitPass                *string
	RPCConnect               *string
	RPCListeners             *[]string
	RPCCert                  *string
	RPCKey                   *string
	RPCMaxClients            *int
	RPCMaxWebsockets         *int
	RPCMaxConcurrentReqs     *int
	RPCQuirks                *bool
	DisableRPC               *bool
	NoTLS                    *bool
	DisableDNSSeed           *bool
	ExternalIPs              *[]string
	Proxy                    *string
	ProxyUser                *string
	ProxyPass                *string
	OnionProxy               *string
	OnionProxyUser           *string
	OnionProxyPass           *string
	Onion                    *bool
	TorIsolation             *bool
	TestNet3                 *bool
	RegressionTest           *bool
	SimNet                   *bool
	AddCheckpoints           *[]string
	DisableCheckpoints       *bool
	DbType                   *string
	Profile                  *int
	CPUProfile               *string
	Upnp                     *bool
	MinRelayTxFee            *float64
	FreeTxRelayLimit         *float64
	NoRelayPriority          *bool
	TrickleInterval          *time.Duration
	MaxOrphanTxs             *int
	Algo                     *string
	Generate                 *bool
	GenThreads               *int
	MiningAddrs              *[]string
	MinerListener            *string
	MinerPass                *string
	BlockMinSize             *int
	BlockMaxSize             *int
	BlockMinWeight           *int
	BlockMaxWeight           *int
	BlockPrioritySize        *int
	UserAgentComments        *[]string
	NoPeerBloomFilters       *bool
	NoCFilters               *bool
	SigCacheMaxSize          *int
	BlocksOnly               *bool
	TxIndex                  *bool
	AddrIndex                *bool
	RelayNonStd              *bool
	RejectNonStd             *bool
	TLSSkipVerify            *bool
	Wallet                   *bool
	NoInitialLoad            *bool
	WalletPass               *string
	WalletServer             *string
	CAFile                   *string
	OneTimeTLSKey            *bool
	ServerTLS                *bool
	LegacyRPCListeners       *[]string
	LegacyRPCMaxClients      *int
	LegacyRPCMaxWebsockets   *int
	ExperimentalRPCListeners *[]string
}
