package main

import (
	"time"

	. "git.parallelcoin.io/dev/9/cmd/config"
	"git.parallelcoin.io/dev/9/cmd/node"
	"git.parallelcoin.io/dev/9/cmd/node/mempool"
)

var NineApp = func() *App {
	return NewApp("9",
		Version("v1.9.9"),
		Tagline("all in one everything for parallelcoin"),
		About("full node, wallet, combined shell, RPC client for the parallelcoin blockchain"),
		DefaultRunner(func(ctx *App) int { return 0 }),
		Cmd("help",
			Pattern("^(h|help)$"),
			Short("show help text and quit"),
			Detail(`	any other command also mentioned with help/h 
	will have its detailed help information printed`),
			Precs("help"),
			Handler(Help),
		),
		Cmd("conf",
			Pattern("^(C|conf)$"),
			Short("run interactive configuration CLI"),
			Detail(`	<datadir> sets the data directory to read and write to`),
			Opts("datadir"),
			Precs("help"),
			Handler(Conf),
		),
		Cmd("new",
			Pattern("^(N|new)$"),
			Short("create new configuration with optional basename and count for testnets"),
			Detail(`	<word> is the basename for the data directories
		<integer> is the number of numbered data directories to create`),
			Opts("word", "integer"),
			Precs("help"),
			Handler(New),
		),
		Cmd("copy",
			Pattern("^(cp|copy)$"),
			Short("create a set of testnet configurations based on a datadir"),
			Detail(`	<datadir> is the base to work from
		<word> is a basename
		<integer> is a number for how many to create`),
			Opts("datadir", "word", "integer"),
			Precs("help"),
			Handler(Copy),
		),
		Cmd("list",
			Pattern("^(l|list|listcommands)$"),
			Short("lists commands available at the RPC endpoint"),
			Detail(`	<datadir> is the enabled data directory
		<ctl> must be present to invoke list
		<wallet> indicates to connect to the wallet RPC
		<node> (or wallet not specified) connect to full node RPC`),
			Opts("datadir", "ctl", "wallet", "node"),
			Precs("help"),
			Handler(List),
		),
		Cmd("ctl",
			Pattern("^(c|ctl)$"),
			Short("sends rpc requests and prints the results"),
			Detail(`	<datadir> sets the data directory to read configurations from
		<node> indicates we are connecting to a full node RPC (overrides wallet and is default)
		<wallet> indicates we are connecting to a wallet RPC
		<word>, <float> and <integer> just cover the items that follow in RPC
		commands the RPC command is expected to be everything after the ctl keyword`),
			Opts("datadir", "node", "wallet", "word", "integer", "float"),
			Precs("help", "list"),
			Handler(Ctl),
		),
		Cmd("node",
			Pattern("^(n|node)$"),
			Short("runs a full node"),
			Detail(`	<datadir> sets the data directory to read configuration and store data`),
			Opts("datadir"),
			Precs("help", "ctl"),
			Handler(Node),
		),
		Cmd("wallet",
			Pattern("^(w|wallet)$"),
			Short("runs a wallet server"),
			Detail(`	<datadir> sets the data directory to read configuration and store data
		<create> runs the wallet create prompt`),
			Opts("datadir", "create"),
			Precs("help", "ctl", "list"),
			Handler(Wallet),
		),
		Cmd("shell",
			Pattern("^(s|shell)$"),
			Short("runs a combined node/wallet server"),
			Detail(`	<datadir> sets the data directory to read configuration and store data
		<create> runs the wallet create prompt`),
			Opts("datadir", "create"),
			Precs("help"),
			Handler(Shell),
		),
		Cmd("mine",
			Pattern("^(m|mine)$"),
			Short("run the standalone miner"),
			Detail(``),
			Opts("datadir"),
			Precs("help"),
			Handler(Mine),
		),
		Cmd("gui",
			Pattern("(^g|gui)$"),
			Short("run the GUI wallet"),
			Detail(``),
			Opts("datadir"),
			Precs("help"),
			Handler(GUI),
		),
		Cmd("test",
			Pattern("^(t|test)$"),
			Short("run multiple full nodes from given <datadir> logging optionally to <datadir>"),
			Detail(`	<datadir> indicates the basename to search for as the path to the test configurations
				<log> indicates to write logs to the individual data directories instead of print to stdout`),
			Opts("word", "log"),
			Precs("help"),
			Handler(TestHandler),
		),
		Cmd("create",
			Pattern("^(cr|create)$"),
			Short("runs the create new wallet prompt"),
			Detail(`	<datadir> sets the data directory where the wallet will be stored`),
			Opts("datadir"),
			Precs("wallet", "shell", "help"),
			Handler(Create),
		),
		Cmd("gencerts",
			Pattern("^(gencerts)$"),
			Short("generate a number of TLS key pairs for nodes"),
			Detail(`	<word> sets the name of the CA signing key file to use
		<integer> sets the number of keys to generate (append a number to the filename)`),
			Opts("word", "integer"),
			Precs("help"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Cmd("genca",
			Pattern("^genca$"),
			Short("generate TLS certification authority key pair"),
			Detail(`	<word> sets the name of the CA signing key file to output`),
			Opts("word"),
			Precs("help"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Cmd("log",
			Pattern("^(L|log)$"),
			Short("write to log in <datadir> file instead of printing to stderr"),
			Detail(`	<datadir> sets the data directory where the wallet will be stored`),
			Opts(),
			Precs("help", "node", "wallet", "shell", "test"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Cmd("datadir",
			Pattern("^(([A-Za-z][:])|[\\~/.]+.*)$"),
			Short("directory to look for configuration or write logs etc"),
			Detail(`	<datadir> sets the data directory where the wallet will be stored`),
			Opts(),
			Precs("help", "node", "ctl", "wallet", "conf", "test", "new", "copy"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Cmd("integer",
			Pattern("^[0-9]+$"),
			Short("number of items to create"),
			Detail(""),
			Opts(),
			Precs("help"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Cmd("float",
			Pattern("^([0-9]+[.][0-9]+)$"),
			Short("a floating point value"),
			Detail(""),
			Opts(),
			Precs("help"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Cmd("word",
			Pattern("^([a-zA-Z0-9][a-zA-Z0-9._-]+)$"),
			Short("mostly used for testnet datadir basenames"),
			Detail(""),
			Opts(),
			Precs("help", "node", "ctl", "wallet", "conf", "test", "new", "copy"),
			Handler(func(args []string, tokens Tokens, app *App) int { return 0 }),
		),
		Group("app",
			Dir("appdatadir",
				Usage("subcommand data directory, sets to datadir/appname if unset"),
			),
			File("cpuprofile",
				Usage("write cpu profile to this file, empty disables cpu profiling"),
			),
			Dir("datadir",
				Default("~/.9"),
				Usage("base folder to keep data for an instance of 9"),
			),
			Dir("logdir",
				Usage("where logs are written, defaults to the appdatadir if unset"),
			),
			Port("profile",
				Usage("http profiling on specified port (1025-65535)"),
			),
			Enable("upnp",
				Usage("enable port forwarding via UPNP"),
			),
		), Group("block",
			Int("maxsize",
				Default(node.DefaultBlockMaxSize),
				Min(node.BlockWeightMin),
				Max(node.BlockSizeMax),
				Usage("max block size in bytes"),
			),
			Int("maxweight",
				Default(node.DefaultBlockMaxWeight),
				Min(node.DefaultBlockMinWeight),
				Max(node.BlockWeightMax),
				Usage("max block weight"),
			),
			Int("minsize",
				Default(node.DefaultBlockMinSize),
				Min(node.DefaultBlockMinSize),
				Max(node.BlockSizeMax),
				Usage("min block size"),
			),
			Int("minweight",
				Default(node.DefaultBlockMinWeight),
				Min(node.DefaultBlockMinWeight),
				Max(node.BlockWeightMax),
				Usage("min block weight"),
			),
			Int("prioritysize",
				Default(mempool.DefaultBlockPrioritySize),
				Min(1000),
				Max(node.BlockSizeMax),
				Usage("the default size for high priority low fee transactions"),
			),
		), Group("chain",
			Tags("addcheckpoints",
				Usage("add checkpoints [height:hash ]*"),
			),
			Enable("disablecheckpoints",
				Usage("disables checkpoints (danger!)"),
			),
			Tag("dbtype",
				Default("ffldb"),
				Usage("set database backend to use for chain"),
			),
			Enabled("addrindex",
				Usage("enable address index (disables also transaction index)"),
			),
			Enabled("txindex",
				Usage("enable transaction index"),
			),
			Enable("rejectnonstd",
				Usage("reject nonstandard transactions even if net parameters allow it"),
			),
			Enable("relaynonstd",
				Usage("relay nonstandard transactions even if net parameters disallow it"),
			),
			Addr("rpc", 11048,
				Default("127.0.0.1:11048"),
				Usage("address of chain rpc to connect to"),
			),
			Int("sigcachemaxsize",
				Default(node.DefaultSigCacheMaxSize),
				Min(1000),
				Max(10000000),
				Usage("max number of signatures to keep in memory"),
			),
		), Group("limit",
			Tag("pass",
				RandomString(32),
				Usage("password for limited user"),
			),
			Tag("user",
				Default("limit"),
				Usage("username with limited privileges"),
			),
		), Group("log",
			Level(
				Default("info"),
				Usage("sets the base default log level"),
			),
			Tags("subsystem",
				Usage("[subsystem:loglevel ]+"),
			),
			Enable("nowrite",
				Usage("disable writing to log file"),
			),
		), Group("mining",
			Tags("addresses",
				Usage("set mining addresses, space separated"),
			),
			Algo("algo",
				Default("random"),
				Usage("select from available mining algorithms"),
			),
			Float("bias",
				Default(-0.5),
				Usage("bias for difficulties -1 = always easy, 1 always hardest"),
			),
			Enable("generate",
				Usage("enable builtin CPU miner"),
			),
			Int("genthreads",
				Default(node.DefaultGenThreads),
				Min(-1),
				Max(4096),
				Usage("set number of threads, -1 = all"),
			),
			Addrs("listener", 11045,
				Usage("set listener address for mining dispatcher"),
			),
			Tag("pass",
				RandomString(32),
				Usage("password to secure mining dispatch connections"),
			),
			Duration("switch",
				Default(time.Second*2),
				Usage("maximum time to mine per round"),
			),
		), Group("p2p",
			Addrs("addpeer", 11047,
				Usage("add permanent p2p peer"),
			),
			Int("banthreshold",
				Default(node.DefaultBanThreshold),
				Usage("how many ban units triggers a ban"),
			),
			Duration("banduration",
				Default(24*time.Hour),
				Usage("how long a ban lasts"),
			),
			Enable("disableban",
				Usage("disables banning peers"),
			),
			Enable("blocksonly",
				Usage("relay only blocks"),
			),
			Addrs("connect", 11047,
				Usage("connect only to these outbound peers"),
			),
			Enable("nolisten",
				Usage("disable p2p listener"),
			),
			Addrs("externalips", 11047,
				Usage("additional external IP addresses to bind to"),
			),
			Float("freetxrelaylimit",
				Default(15.0),
				Usage("limit of 'free' relay in thousand bytes per minute"),
			),
			Addrs("listen", 11047,
				Default("127.0.0.1:11047"),
				Usage("addresss to listen on for p2p connections"),
			),
			Int("maxorphantxs",
				Default(node.DefaultMaxOrphanTransactions),
				Min(0),
				Max(10000),
				Usage("maximum number of orphan transactions to keep in memory"),
			),
			Int("maxpeers",
				Default(node.DefaultMaxPeers),
				Min(2),
				Max(1024),
				Usage("maximum number of peers to connect to"),
			),
			Float("minrelaytxfee",
				Default(0.0001),
				Usage("minimum relay tx fee, baseline considered to be zero for relay"),
			),
			Net("network",
				Default("mainnet"),
				Usage("network to connect to"),
			),
			Enable("nobanning",
				Usage("disable banning of peers"),
			),
			Enable("nobloomfilters",
				Usage("disable bloom filters"),
			),
			Enable("nocfilters",
				Usage("disable cfilters"),
			),
			Enable("nodns",
				Usage("disable DNS seeding"),
			),
			Enable("norelaypriority",
				Usage("disables prioritisation of relayed transactions"),
			),
			Duration("trickleinterval",
				Default(time.Second*27),
				Usage("minimum time between attempts to send new inventory to a connected peer"),
			),
			Tags("useragentcomments",
				Usage("comment to add to version identifier for node"),
			),
			Addrs("whitelist", 11047,
				Usage("peers who are never banned"),
			),
		),
		Group("proxy",
			Addr("address", 9050,
				Usage("address of socks proxy"),
			),
			Enable("isolation",
				Usage("enable randomisation of tor login to separate streams"),
			),
			Tag("pass",
				RandomString(32),
				Usage("password for proxy"),
			),
			Enable("tor",
				Usage("proxy is a tor proxy"),
			),
			Tag("user",
				Default("user"),
				Usage("username for proxy"),
			),
		),
		Group("rpc",
			Addr("connect", 11048,
				Default("127.0.0.1:11048"),
				Usage("connect to this node RPC endpoint"),
			),
			Enable("disable",
				Usage("disable rpc server"),
			),
			Addrs("listen", 11048,
				Default("127.0.0.1:11048"),
				Usage("address to listen for node rpc clients"),
			),
			Int("maxclients",
				Default(node.DefaultMaxRPCClients),
				Min(2),
				Max(1024),
				Usage("max clients for rpc"),
			),
			Int("maxconcurrentreqs",
				Default(node.DefaultMaxRPCConcurrentReqs),
				Min(2),
				Max(1024),
				Usage("maximum concurrent requests to handle"),
			),
			Int("maxwebsockets",
				Default(node.DefaultMaxRPCWebsockets),
				Max(1024),
				Usage("maximum websockets clients"),
			),
			Tag("pass",
				RandomString(32),
				Usage("password for rpc services"),
			),
			Enable("quirks",
				Usage("enable json rpc quirks matching bitcoin core"),
			),
			Tag("user",
				Default("user"),
				Usage("username for rpc services"),
			),
		),
		Group("tls",
			File("key",
				Default("tls.key"),
				Usage("file containing tls key"),
			),
			File("cert",
				Default("tls.cert"),
				Usage("file containing tls certificate"),
			),
			File("cafile",
				Default("tls.cafile"),
				Usage("set the certificate authority file to use for verifying rpc connections"),
			),
			Enable("disable",
				Usage("disable SSL on RPC connections"),
			),
			Enable("onetime",
				Usage("creates a key pair but does not write the secret for future runs"),
			),
			Enabled("server",
				Usage("enable tls for RPC servers"),
			),
			Enable("skipverify",
				Usage("skip verifying tls certificates with CAFile"),
			),
		),
		Group("wallet",
			Addr("server", 11046,
				Default("127.0.0.1:11046"),
				Usage("address of wallet rpc to connect to"),
			),
			Enable("noinitialload",
				Usage("disable automatic opening of the wallet at startup"),
			),
			Tag("pass",
				RandomString(32),
				Usage("password for the non-own transaction data in the wallet"),
			),
			Enable("enable",
				Usage("use configured wallet rpc instead of full node"),
			),
		),
	)
}
