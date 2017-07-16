package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tsavola/wag/wasm"
	"golang.org/x/crypto/acme/autocert"

	"github.com/tsavola/gate/run"
	"github.com/tsavola/gate/server"
	_ "github.com/tsavola/gate/service/defaults"
	"github.com/tsavola/gate/service/origin"
)

const (
	renewCertBefore = 30 * 24 * time.Hour

	memorySizeLimit = 256 * wasm.Page
	stackSize       = 16 * 4096
)

func main() {
	var (
		config = run.Config{
			LibDir:   "lib",
			MaxProcs: run.DefaultMaxProcs,
		}
		addr         = "localhost:8888"
		letsencrypt  = false
		email        = ""
		acceptTOS    = false
		certCacheDir = "/var/lib/gate-httpserver-letsencrypt"
		debug        = false
	)

	flag.StringVar(&config.LibDir, "libdir", config.LibDir, "path")
	flag.UintVar(&config.Uids[0], "boot-uid", config.Uids[0], "user id for bootstrapping executor")
	flag.UintVar(&config.Gids[0], "boot-gid", config.Gids[0], "group id for bootstrapping executor")
	flag.UintVar(&config.Uids[1], "exec-uid", config.Uids[1], "user id for executing code")
	flag.UintVar(&config.Gids[1], "exec-gid", config.Gids[1], "group id for executing code")
	flag.UintVar(&config.Gids[2], "pipe-gid", config.Gids[2], "group id for file descriptor sharing")
	flag.IntVar(&config.MaxProcs, "max-procs", config.MaxProcs, "limit number of simultaneous programs")
	flag.StringVar(&addr, "addr", addr, "listening [address]:port")
	flag.BoolVar(&letsencrypt, "letsencrypt", letsencrypt, "enable automatic TLS; domain names should be listed after the options")
	flag.StringVar(&email, "email", email, "contact address for Let's Encrypt")
	flag.BoolVar(&acceptTOS, "accept-tos", acceptTOS, "accept Let's Encrypt's terms of service")
	flag.StringVar(&certCacheDir, "cert-cache-dir", certCacheDir, "certificate storage")
	flag.BoolVar(&debug, "debug", debug, "write payload programs' debug output to stderr")
	flag.Parse()
	domains := flag.Args()

	env, err := run.NewEnvironment(&config)
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	e := server.Executor{
		MemorySizeLimit: memorySizeLimit,
		StackSize:       stackSize,
		Env:             env,
		Services:        services,
		Log:             log.New(os.Stderr, "", 0),
	}

	if debug {
		e.Debug = os.Stderr
	}

	http.Handle("/", e.Handler())

	if letsencrypt {
		if !acceptTOS {
			log.Fatal("-accept-tos option not set")
		}

		m := autocert.Manager{
			Prompt:      autocert.AcceptTOS,
			Cache:       autocert.DirCache(certCacheDir),
			HostPolicy:  autocert.HostWhitelist(domains...),
			RenewBefore: renewCertBefore,
			Email:       email,
		}

		s := http.Server{
			Addr: addr,
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}

		err = s.ListenAndServeTLS("", "")
	} else {
		err = http.ListenAndServe(addr, nil)
	}

	log.Fatal(err)
}

func services(r io.Reader, w io.Writer) run.ServiceRegistry {
	return origin.CloneRegistryWith(nil, r, w)
}
