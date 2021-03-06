package node
import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"
	"git.parallelcoin.io/dev/9/cmd/nine"
	indexers "git.parallelcoin.io/dev/9/pkg/chain/index"
	database "git.parallelcoin.io/dev/9/pkg/db"
	cl "git.parallelcoin.io/dev/9/pkg/util/cl"
	"git.parallelcoin.io/dev/9/pkg/util/interrupt"
)
// blockDbNamePrefix is the prefix for the block database name.  The database type is appended to this value to form the full block database name.
const blockDbNamePrefix = "blocks"
var StateCfg = &nine.StateConfig{}
var Cfg = &nine.Config{}
// // winServiceMain is only invoked on Windows.  It detects when pod is running as a service and reacts accordingly.
// var winServiceMain func() (bool, error)
// Main is the real main function for pod.  It is necessary to work around the fact that deferred functions do not run when os.Exit() is called.  The optional serverChan parameter is mainly used by the service code to be notified with the server once it is setup so it can gracefully stop it when requested from the service control manager.
func Main(serverChan chan<- *server, started chan struct{}) (err error) {
	// // Call serviceMain on Windows to handle running as a service.  When
	// // the return isService flag is true, exit now since we ran as a
	// // service.  Otherwise, just fall through to normal operation.
	// if runtime.GOOS == "windows" {
	// 	isService, err := winServiceMain()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	if isService {
	// 		os.Exit(0)
	// 	}
	// }
	shutdownChan := make(chan struct{})
	interrupt.AddHandler(
		func() {
			log <- cl.Inf("closing shutdown channel")
			close(shutdownChan)
		},
	)
	// Show version at startup.
	log <- cl.Info{"version", Version()}
	// Enable http profiling server if requested.
	if Cfg.Profile != nil {
		log <- cl.Dbg("profiling requested")
		go func() {
			listenAddr :=
				net.JoinHostPort("", fmt.Sprint(*Cfg.Profile))
			log <- cl.Info{"profile server listening on", listenAddr}
			profileRedirect :=
				http.RedirectHandler("/debug/pprof",
					http.StatusSeeOther)
			http.Handle("/", profileRedirect)
			log <- cl.Error{"profile server", http.ListenAndServe(listenAddr, nil)}
		}()
	}
	// Write cpu profile if requested.
	if Cfg.CPUProfile != nil {
		var f *os.File
		f, err = os.Create(*Cfg.CPUProfile)
		if err != nil {
			log <- cl.Error{"unable to create cpu profile:", err}
			return
		}
		e := pprof.StartCPUProfile(f)
		if e != nil {
			log <- cl.Warn{"failed to start up cpu profiler:", e}
		}
		defer f.Close()
		defer pprof.StopCPUProfile()
	}
	// Perform upgrades to pod as new versions require it.
	if err = doUpgrades(); err != nil {
		log <- cl.Error{err}
		return
	}
	// Return now if an interrupt signal was triggered.
	if interrupt.Requested() {
		return nil
	}
	// Load the block database.
	var db database.DB
	log <- cl.Debug{
		"loading db with", ActiveNetParams.Params.Name}
	db, err = loadBlockDB()
	if err != nil {
		log <- cl.Error{err}
		return
	}
	defer func() {
		// Ensure the database is sync'd and closed on shutdown.
		log <- cl.Inf("gracefully shutting down the database...")
		db.Close()
	}()
	// Return now if an interrupt signal was triggered.
	if interrupt.Requested() {
		return nil
	}
	// Drop indexes and exit if requested. NOTE: The order is important here because dropping the tx index also drops the address index since it relies on it.
	if StateCfg.DropAddrIndex {
		log <- cl.Warn{"dropping address index"}
		if err = indexers.DropAddrIndex(db, interrupt.ShutdownRequestChan); err != nil {
			log <- cl.Error{err}
			if err != nil {
				return
			}
		}
	}
	if StateCfg.DropTxIndex {
		log <- cl.Warn{"dropping transaction index"}
		if err = indexers.DropTxIndex(db, interrupt.ShutdownRequestChan); err != nil {
			log <- cl.Error{err}
			if err != nil {
				return
			}
		}
	}
	if StateCfg.DropCfIndex {
		log <- cl.Warn{"dropping cfilter index"}
		if err = indexers.DropCfIndex(db, interrupt.ShutdownRequestChan); err != nil {
			log <- cl.Error{err}
			if err != nil {
				return
			}
		}
	}
	// Create server and start it.
	server, err := newServer(*Cfg.Listeners, db, ActiveNetParams.Params, interrupt.ShutdownRequestChan, *Cfg.Algo)
	if err != nil {
		log <- cl.Errorf{
			"unable to start server on %v: %v", *Cfg.Listeners, err}
		return err
	}
	interrupt.AddHandler(
		func() {
			log <- cl.Inf("gracefully shutting down the server...")
			e := server.Stop()
			if e != nil {
				log <- cl.Warn{"failed to stop server", e}
			}
			server.WaitForShutdown()
			log <- cl.Inf("server shutdown complete")
		},
	)
	server.Start()
	if serverChan != nil {
		serverChan <- server
	}
	close(started)
	log <- cl.Info{"blockchain node is now started"}
	// Wait until the interrupt signal is received from an OS signal or shutdown is requested through one of the subsystems such as the RPC server.
	<-interrupt.HandlersDone
	return nil
}
// dbPath returns the path to the block database given a database type.
func blockDbPath(dbType string) string {
	// The database name is based on the database type.
	dbName := blockDbNamePrefix + "_" + dbType
	if dbType == "sqlite" {
		dbName += ".db"
	}
	dbPath := filepath.Join(
		filepath.Join(
			*Cfg.AppDataDir, NetName(ActiveNetParams)), dbName)
	return dbPath
}
// loadBlockDB loads (or creates when needed) the block database taking into account the selected database backend and returns a handle to it.  It also additional logic such warning the user if there are multiple databases which consume space on the file system and ensuring the regression test database is clean when in regression test mode.
func loadBlockDB() (database.DB, error) {
	// The memdb backend does not have a file path associated with it, so handle it uniquely.  We also don't want to worry about the multiple database type warnings when running with the memory database.
	if *Cfg.DbType == "memdb" {
		log <- cl.Inf("creating block database in memory")
		db, err := database.Create(*Cfg.DbType)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	warnMultipleDBs()
	// The database name is based on the database type.
	dbPath := blockDbPath(*Cfg.DbType)
	// The regression test is special in that it needs a clean database for each run, so remove it now if it already exists.
	e := removeRegressionDB(dbPath)
	if e != nil {
		log <- cl.Debug{"failed to remove regression db:", e}
	}
	log <- cl.Infof{"loading block database from '%s'", dbPath}
	db, err := database.Open(*Cfg.DbType, dbPath, ActiveNetParams.Net)
	if err != nil {
		// Return the error if it's not because the database doesn't exist.
		if dbErr, ok := err.(database.Error); !ok || dbErr.ErrorCode !=
			database.ErrDbDoesNotExist {
			return nil, err
		}
		// Create the db if it does not exist.
		err = os.MkdirAll(*Cfg.DataDir, 0700)
		if err != nil {
			return nil, err
		}
		db, err = database.Create(*Cfg.DbType, dbPath, ActiveNetParams.Net)
		if err != nil {
			return nil, err
		}
	}
	log <- cl.Inf("block database loaded")
	return db, nil
}
/*
func PreMain() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Block and transaction processing can cause bursty allocations.  This limits the garbage collector from excessively overallocating during bursts.  This value was arrived at with the help of profiling live usage.
	debug.SetGCPercent(10)
	// Up some limits.
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	// Call serviceMain on Windows to handle running as a service.  When the return isService flag is true, exit now since we ran as a service.  Otherwise, just fall through to normal operation.
	if runtime.GOOS == "windows" {
		isService, err := winServiceMain()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if isService {
			os.Exit(0)
		}
	}
	// Work around defer not working after os.Exit()
	if err := Main(nil); err != nil {
		os.Exit(1)
	}
}
*/
// removeRegressionDB removes the existing regression test database if running in regression test mode and it already exists.
func removeRegressionDB(
	dbPath string,
) error {
	// Don't do anything if not in regression test mode.
	if !*Cfg.RegressionTest {
		log <- cl.Debug{"not in regression mode"}
		return nil
	}
	// Remove the old regression test database if it already exists.
	fi, err := os.Stat(dbPath)
	if err == nil {
		log <- cl.Infof{"removing regression test database from '%s'", dbPath}
		if fi.IsDir() {
			err := os.RemoveAll(dbPath)
			if err != nil {
				return err
			}
		} else {
			err := os.Remove(dbPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
// warnMultipleDBs shows a warning if multiple block database types are detected. This is not a situation most users want.  It is handy for development however to support multiple side-by-side databases.
func warnMultipleDBs() {
	// This is intentionally not using the known db types which depend on the database types compiled into the binary since we want to detect legacy db types as well.
	dbTypes := []string{"ffldb", "leveldb", "sqlite"}
	duplicateDbPaths := make([]string, 0, len(dbTypes)-1)
	for _, dbType := range dbTypes {
		if dbType == *Cfg.DbType {
			continue
		}
		// Store db path as a duplicate db if it exists.
		dbPath := blockDbPath(dbType)
		if FileExists(dbPath) {
			duplicateDbPaths = append(duplicateDbPaths, dbPath)
		}
	}
	// Warn if there are extra databases.
	if len(duplicateDbPaths) > 0 {
		selectedDbPath := blockDbPath(*Cfg.DbType)
		log <- cl.Warnf{
			"\nThere are multiple block chain databases using different database types.\n" +
				"You probably don't want to waste disk space by having more than one.\n" +
				"Your current database is located at [%v].\n" +
				"The additional database is located at %v",
			selectedDbPath,
			duplicateDbPaths,
		}
	}
}
