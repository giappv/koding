package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"koding/artifact"
	"koding/db/mongodb/modelhelper"
	"koding/kites/kloud/keys"
	"koding/kites/kloud/koding"
	"koding/kites/kloud/multiec2"

	"koding/kites/kloud/klient"
	"koding/kites/kloud/kloud"
	kloudprotocol "koding/kites/kloud/protocol"

	"github.com/koding/metrics"

	"github.com/koding/kite"
	kiteconfig "github.com/koding/kite/config"
	"github.com/koding/kite/protocol"
	"github.com/koding/logging"
	"github.com/koding/multiconfig"
	"github.com/mitchellh/goamz/aws"
)

var Name = "kloud"

// Config defines the configuration that Kloud needs to operate.
type Config struct {
	// ---  KLOUD SPECIFIC ---
	IP          string
	Port        int
	Region      string
	Environment string

	// Connect to Koding mongodb
	MongoURL string `required:"true"`

	// Endpoint for fetchin plans
	PlanEndpoint string `required:"true"`

	// --- DEVELOPMENT CONFIG ---
	// Show version and exit if enabled
	Version bool

	// Enable debug log mode
	DebugMode bool

	// Enable production mode, operates on production channel
	ProdMode bool

	// Enable test mode, disabled some authentication checks
	TestMode bool

	// Defines the base domain for domain creation
	HostedZone string `required:"true"`

	// Defines the default AMI Tag to use for koding provider
	AMITag string

	// --- KLIENT DEVELOPMENT ---
	// KontrolURL to connect and to de deployed with klient
	KontrolURL string `required:"true"`

	// Private key to create kite.key
	PrivateKey string `required:"true"`

	// Public key to create kite.key
	PublicKey string `required:"true"`

	// --- KONTROL CONFIGURATION ---
	Public      bool   // Try to register with a public ip
	Proxy       bool   // Try to register behind a koding proxy
	RegisterURL string // Explicitly register with this given url

	// Artifacts endpoint port
	ArtifactPort int
}

func main() {
	conf := new(Config)

	// Load the config, it's reads environment variables or from flags
	multiconfig.New().MustLoad(conf)

	if conf.Version {
		fmt.Println(kloud.VERSION)
		os.Exit(0)
	}

	k := newKite(conf)

	if conf.DebugMode {
		k.Log.Info("Debug mode enabled")
	}

	if conf.TestMode {
		k.Log.Info("Test mode enabled")
	}

	registerURL := k.RegisterURL(!conf.Public)
	if conf.RegisterURL != "" {
		u, err := url.Parse(conf.RegisterURL)
		if err != nil {
			k.Log.Fatal("Couldn't parse register url: %s", err)
		}

		registerURL = u
	}

	if conf.Proxy {
		k.Log.Info("Proxy mode is enabled")
		// Koding proxies in production only
		proxyQuery := &protocol.KontrolQuery{
			Username:    "koding",
			Environment: "production",
			Name:        "proxy",
		}

		k.Log.Info("Seaching proxy: %#v", proxyQuery)
		go k.RegisterToProxy(registerURL, proxyQuery)
	} else {
		if err := k.RegisterForever(registerURL); err != nil {
			k.Log.Fatal(err.Error())
		}
	}

	// TODO use kite's http server instead of creating another one here
	// this is used for application lifecycle management
	go artifact.StartDefaultServer(Name, conf.ArtifactPort)

	k.Run()
}

func newKite(conf *Config) *kite.Kite {
	k := kite.New(kloud.NAME, kloud.VERSION)
	k.Config = kiteconfig.MustGet()
	k.Config.Port = conf.Port

	if conf.Region != "" {
		k.Config.Region = conf.Region
	}

	if conf.Environment != "" {
		k.Config.Environment = conf.Environment
	}

	if conf.AMITag != "" {
		k.Log.Warning("Default AMI Tag changed from %s to %s", koding.DefaultCustomAMITag, conf.AMITag)
		koding.DefaultCustomAMITag = conf.AMITag
	}

	klientFolder := "development/latest"
	checkInterval := time.Second * 5
	if conf.ProdMode {
		k.Log.Info("Prod mode enabled")
		klientFolder = "production/latest"
		checkInterval = time.Millisecond * 500
	}
	k.Log.Info("Klient distribution channel is: %s", klientFolder)

	modelhelper.Initialize(conf.MongoURL)
	db := modelhelper.Mongo

	kontrolPrivateKey, kontrolPublicKey := kontrolKeys(conf)

	// Credential belongs to the `koding-kloud` user in AWS IAM's
	auth := aws.Auth{
		AccessKey: "AKIAIKAVWAYVSMCW4Z5A",
		SecretKey: "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",
	}

	stats, err := metrics.NewDogStatsD("kloud.aws")
	if err != nil {
		panic(err)
	}

	dnsInstance := koding.NewDNSClient(conf.HostedZone, auth)
	domainStorage := koding.NewDomainStorage(db)

	kodingProvider := &koding.Provider{
		Kite:              k,
		Log:               newLogger("koding", conf.DebugMode),
		Session:           db,
		DomainStorage:     domainStorage,
		EC2Clients:        multiec2.New(auth, []string{"us-east-1", "ap-southeast-1"}),
		DNS:               dnsInstance,
		Bucket:            koding.NewBucket("koding-klient", klientFolder, auth),
		Test:              conf.TestMode,
		KontrolURL:        getKontrolURL(conf.KontrolURL),
		KontrolPrivateKey: kontrolPrivateKey,
		KontrolPublicKey:  kontrolPublicKey,
		KeyName:           keys.DeployKeyName,
		PublicKey:         keys.DeployPublicKey,
		PrivateKey:        keys.DeployPrivateKey,
		KlientPool:        klient.NewPool(k),
		InactiveMachines:  make(map[string]*time.Timer),
		Stats:             stats,
	}

	// be sure it satisfies the provider interface
	var _ kloudprotocol.Provider = kodingProvider

	kodingProvider.PlanChecker = func(m *kloudprotocol.Machine) (koding.Checker, error) {
		a, err := kodingProvider.NewClient(m)
		if err != nil {
			return nil, err
		}

		return &koding.PlanChecker{
			Api:      a,
			Provider: kodingProvider,
			DB:       kodingProvider.Session,
			Kite:     kodingProvider.Kite,
			Log:      kodingProvider.Log,
			Username: m.Username,
			Machine:  m,
		}, nil
	}

	kodingProvider.PlanFetcher = func(m *kloudprotocol.Machine) (koding.Plan, error) {
		return kodingProvider.Fetcher(conf.PlanEndpoint, m)
	}

	go kodingProvider.RunChecker(checkInterval)
	go kodingProvider.RunCleaners(time.Minute)

	kld := kloud.NewWithDefaults()
	kld.Storage = kodingProvider
	kld.DomainStorage = domainStorage
	kld.Domainer = dnsInstance
	kld.Locker = kodingProvider
	kld.Log = newLogger(Name, conf.DebugMode)

	err = kld.AddProvider("koding", kodingProvider)
	if err != nil {
		panic(err)
	}

	// Admin bypass if the username is koding or kloud
	k.PreHandleFunc(func(r *kite.Request) (interface{}, error) {
		if r.Args == nil {
			return nil, nil
		}

		if _, err := r.Args.SliceOfLength(1); err != nil {
			return nil, nil
		}

		var args struct {
			MachineId string
			Username  string
		}

		if err := r.Args.One().Unmarshal(&args); err != nil {
			return nil, nil
		}

		if koding.IsAdmin(r.Username) && args.Username != "" {
			k.Log.Warning("[%s] ADMIN COMMAND: replacing username from '%s' to '%s'",
				args.MachineId, r.Username, args.Username)
			r.Username = args.Username
		}

		return nil, nil
	})

	// Machine handling methods
	k.HandleFunc("build", kld.Build)
	k.HandleFunc("start", kld.Start)
	k.HandleFunc("stop", kld.Stop)
	k.HandleFunc("restart", kld.Restart)
	k.HandleFunc("info", kld.Info)
	k.HandleFunc("destroy", kld.Destroy)
	k.HandleFunc("event", kld.Event)
	k.HandleFunc("resize", kld.Resize)
	k.HandleFunc("reinit", kld.Reinit)

	// Domain records handling methods
	k.HandleFunc("domain.set", kld.DomainSet)
	k.HandleFunc("domain.unset", kld.DomainUnset)
	k.HandleFunc("domain.add", kld.DomainAdd)
	k.HandleFunc("domain.remove", kld.DomainRemove)

	return k
}

func newLogger(name string, debug bool) logging.Logger {
	log := logging.NewLogger(name)
	logHandler := logging.NewWriterHandler(os.Stderr)
	logHandler.Colorize = true
	log.SetHandler(logHandler)

	if debug {
		log.SetLevel(logging.DEBUG)
		logHandler.SetLevel(logging.DEBUG)
	}

	return log
}

func kontrolKeys(conf *Config) (string, string) {
	pubKey, err := ioutil.ReadFile(conf.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}
	publicKey := string(pubKey)

	privKey, err := ioutil.ReadFile(conf.PrivateKey)
	if err != nil {
		log.Fatalln(err)
	}
	privateKey := string(privKey)

	return privateKey, publicKey
}

func getKontrolURL(ownURL string) string {
	// read kontrolURL from kite.key if it doesn't exist.
	kontrolURL := kiteconfig.MustGet().KontrolURL

	if ownURL != "" {
		u, err := url.Parse(ownURL)
		if err != nil {
			log.Fatalln(err)
		}

		kontrolURL = u.String()
	}

	return kontrolURL
}
