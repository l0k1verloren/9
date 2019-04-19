package main

import (
	"git.parallelcoin.io/dev/9/cmd/config"
	"git.parallelcoin.io/dev/9/cmd/nine"
)

func MakeConfig(c *config.App) (out *nine.Config) {
	C := c.Cats
	var configFile string
	var tn, sn, rn bool
	out = &nine.Config{
		ConfigFile:               &configFile,
		AppDataDir:               C.Str("app", "appdatadir"),
		DataDir:                  C.Str("app", "datadir"),
		LogDir:                   C.Str("app", "logdir"),
		LogLevel:                 C.Str("log", "level"),
		Subsystems:               C.Map("log", "subsystem"),
		Network:                  C.Str("p2p", "network"),
		AddPeers:                 C.Tags("p2p", "addpeer"),
		ConnectPeers:             C.Tags("p2p", "connect"),
		MaxPeers:                 C.Int("p2p", "maxpeers"),
		Listeners:                C.Tags("p2p", "listen"),
		DisableListen:            C.Bool("p2p", "nolisten"),
		DisableBanning:           C.Bool("p2p", "disableban"),
		BanDuration:              C.Duration("p2p", "banduration"),
		BanThreshold:             C.Int("p2p", "banthreshold"),
		Whitelists:               C.Tags("p2p", "whitelist"),
		Username:                 C.Str("rpc", "user"),
		Password:                 C.Str("rpc", "pass"),
		ServerUser:               C.Str("rpc", "user"),
		ServerPass:               C.Str("rpc", "pass"),
		LimitUser:                C.Str("limit", "user"),
		LimitPass:                C.Str("limit", "pass"),
		RPCConnect:               C.Str("rpc", "connect"),
		RPCListeners:             C.Tags("rpc", "listen"),
		RPCCert:                  C.Str("tls", "cert"),
		RPCKey:                   C.Str("tls", "key"),
		RPCMaxClients:            C.Int("rpc", "maxclients"),
		RPCMaxWebsockets:         C.Int("rpc", "maxwebsockets"),
		RPCMaxConcurrentReqs:     C.Int("rpc", "maxconcurrentreqs"),
		RPCQuirks:                C.Bool("rpc", "quirks"),
		DisableRPC:               C.Bool("rpc", "disable"),
		NoTLS:                    C.Bool("tls", "disable"),
		DisableDNSSeed:           C.Bool("p2p", "nodns"),
		ExternalIPs:              C.Tags("p2p", "externalips"),
		Proxy:                    C.Str("proxy", "address"),
		ProxyUser:                C.Str("proxy", "user"),
		ProxyPass:                C.Str("proxy", "pass"),
		OnionProxy:               C.Str("proxy", "address"),
		OnionProxyUser:           C.Str("proxy", "user"),
		OnionProxyPass:           C.Str("proxy", "pass"),
		Onion:                    C.Bool("proxy", "tor"),
		TorIsolation:             C.Bool("proxy", "isolation"),
		TestNet3:                 &tn,
		RegressionTest:           &rn,
		SimNet:                   &sn,
		AddCheckpoints:           C.Tags("chain", "addcheckpoints"),
		DisableCheckpoints:       C.Bool("chain", "disablecheckpoints"),
		DbType:                   C.Str("chain", "dbtype"),
		Profile:                  C.Int("app", "profile"),
		CPUProfile:               C.Str("app", "cpuprofile"),
		Upnp:                     C.Bool("app", "upnp"),
		MinRelayTxFee:            C.Float("p2p", "minrelaytxfee"),
		FreeTxRelayLimit:         C.Float("p2p", "freetxrelaylimit"),
		NoRelayPriority:          C.Bool("p2p", "norelaypriority"),
		TrickleInterval:          C.Duration("p2p", "trickleinterval"),
		MaxOrphanTxs:             C.Int("p2p", "maxorphantxs"),
		Algo:                     C.Str("mining", "algo"),
		Generate:                 C.Bool("mining", "generate"),
		GenThreads:               C.Int("mining", "genthreads"),
		MiningAddrs:              C.Tags("mining", "addresses"),
		MinerListener:            C.Str("mining", "listener"),
		MinerPass:                C.Str("mining", "pass"),
		BlockMinSize:             C.Int("block", "minsize"),
		BlockMaxSize:             C.Int("block", "maxsize"),
		BlockMinWeight:           C.Int("block", "minweight"),
		BlockMaxWeight:           C.Int("block", "maxweight"),
		BlockPrioritySize:        C.Int("block", "prioritysize"),
		UserAgentComments:        C.Tags("p2p", "useragentcomments"),
		NoPeerBloomFilters:       C.Bool("p2p", "nobloomfilters"),
		NoCFilters:               C.Bool("p2p", "nocfilters"),
		SigCacheMaxSize:          C.Int("chain", "sigcachemaxsize"),
		BlocksOnly:               C.Bool("p2p", "blocksonly"),
		TxIndex:                  C.Bool("chain", "txindex"),
		AddrIndex:                C.Bool("chain", "addrindex"),
		RelayNonStd:              C.Bool("chain", "relaynonstd"),
		RejectNonStd:             C.Bool("chain", "rejectnonstd"),
		TLSSkipVerify:            C.Bool("tls", "skipverify"),
		Wallet:                   C.Bool("wallet", "enable"),
		NoInitialLoad:            C.Bool("wallet", "noinitialload"),
		WalletPass:               C.Str("wallet", "pass"),
		WalletServer:             C.Str("rpc", "wallet"),
		CAFile:                   C.Str("tls", "cafile"),
		OneTimeTLSKey:            C.Bool("tls", "onetime"),
		ServerTLS:                C.Bool("tls", "server"),
		LegacyRPCListeners:       C.Tags("rpc", "listen"),
		LegacyRPCMaxClients:      C.Int("rpc", "maxclients"),
		LegacyRPCMaxWebsockets:   C.Int("rpc", "maxwebsockets"),
		ExperimentalRPCListeners: &[]string{},
	}
	return
}
