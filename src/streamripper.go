package main

import "crypto/rand"
import "crypto/tls"
import "crypto/x509"
import "encoding/binary"
import "flag"
//import "fmt"
import "io/ioutil"
import rnd "math/rand"
import "log"
import "net"
import "os"
import "sync"
import "time"

import "ripper"

type RipperConfig struct {
	RunServer *bool
	RunClient *bool

	KeyPath *string
	CertPath *string

	ClientDelay *int

	BlockSize *int
	Threads *int

	Source *string
	Destination *string
}


type ConnectionData struct {
	active bool
	connectionId uint64
	inbound bool
	j ripper.Joiner
	s ripper.Splitter
	closeOnce sync.Once 
}

var ConnectionList map[uint64]*ConnectionData;

func handleServerConnection(conn net.Conn, config RipperConfig) {
	var connectionId uint64
	var err = binary.Read(conn, binary.LittleEndian, &connectionId)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	log.Printf("ConnectionId %v connected.", connectionId)

	var connData *ConnectionData;

	connData, ok := ConnectionList[connectionId];

	if ok == false {
		connData = &ConnectionData{};
		connData.active = true;
		connData.connectionId = connectionId;
		connData.inbound = true;
		ConnectionList[connectionId] = connData;
	
		connData.j = ripper.Joiner{}
		connData.j.Blocksize = uint32(*config.BlockSize)
		connData.j.BlocksIn = make(chan ripper.SplitBlock, 100)
		connData.j.BlocksOut = make(chan ripper.RawBlock, 100)
	
		go connData.j.RunJoiner()

		connData.j.AddOut(os.Stdout)

		go func() {
			connData.j.WaitOut.Wait()
			log.Printf("Finished stdout")

			os.Exit(0)
		}()

	}

	go func() {
		connData.j.AddIn(conn)
		connData.j.WaitIn.Wait()
		connData.closeOnce.Do(func() {
			close(connData.j.BlocksOut)
		})
		log.Printf("Connection from %v closed.", conn.RemoteAddr())
		conn.Close()
	}()
}

func server(config RipperConfig) {
	certBytes, err := ioutil.ReadFile(*config.CertPath)

	if err != nil {
		log.Fatalln("Unable to read cert", err)
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Unable to add certificate to certificate pool")
	}

	cert, err := tls.LoadX509KeyPair(*config.CertPath, *config.KeyPath)

	if err != nil {
		log.Fatal(err)
	}

	tlsconfig := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Ensure that we only use our "CA" to validate certificates
		ClientCAs: clientCertPool,
		// TLS 1.2 because we can
		MinVersion: tls.VersionTLS12,
	}

	tlsconfig.Rand = rand.Reader
	tlsconfig.BuildNameToCertificate()

	ln, err := tls.Listen("tcp", *config.Source, &tlsconfig)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// The handle server connection will accept one new connection at a time.
		// This allows multiple connections to be started at the same time, but allows 
		// us to handle starting them one at a time.
		handleServerConnection(conn, config);
	}
}

func startClientThread(config RipperConfig, connData *ConnectionData) {
	certBytes, err := ioutil.ReadFile(*config.CertPath)

	if err != nil {
		log.Fatalln("Unable to read cert", err)
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Unable to add certificate to certificate pool")
	}

	cert, err := tls.LoadX509KeyPair(*config.CertPath, *config.KeyPath)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: clientCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	conn, err := tls.Dial("tcp", *config.Destination, tlsConfig)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer conn.Close()

		err = binary.Write(conn, binary.LittleEndian, connData.connectionId)

		if err != nil {
			log.Fatal(err)
		}

		connData.s.AddOut(conn);
		connData.s.WaitOut.Wait();
	}()
}

func client(config RipperConfig) {


	var connData *ConnectionData;

	connData = &ConnectionData{};
	connData.active = true;
	var connectionId uint64 = rnd.Uint64()
	connData.connectionId = connectionId;
	connData.inbound = false;

	connData.s = ripper.Splitter{}
	connData.s.Blocksize = uint32(*config.BlockSize)
	connData.s.BlocksIn = make(chan ripper.RawBlock, 100)
	connData.s.BlocksOut = make(chan ripper.SplitBlock, 100)

	go connData.s.RunSplitter()

	for i := 0; i < *config.Threads; i++ {
		startClientThread(config, connData)
	}

	connData.s.AddIn(os.Stdin)

	connData.s.WaitIn.Wait()
	connData.s.WaitOut.Wait()
	log.Fatal("Done!")
}

func main() {
	rnd.Seed(time.Now().UTC().UnixNano())

	ConnectionList = make(map[uint64]*ConnectionData);

	var config RipperConfig

	config.RunServer = flag.Bool("server", false, "Run as Server")
	config.RunClient = flag.Bool("client", false, "Run as Client")

	config.KeyPath = flag.String("key", "mykey.pem", "Key")
	config.CertPath = flag.String("cert", "mycert.pem", "Cert")

	config.ClientDelay = flag.Int("clientdelay", 0, "Delay in seconds before connecting to server")

	config.BlockSize = flag.Int("buffersize", 1048576, "Buffer size for each block")
	config.Threads = flag.Int("threads", 4, "Number of threads")
	config.Source = flag.String("src", "localhost:8887", "Source ( - stdin )")
	config.Destination = flag.String("dest", "localhost:8888", "Destination ( - stdout )")

	flag.Parse()

	if *config.RunServer {
		if *config.Source == "-" {
			log.Fatalln("Server cannot read from stdin")
		}
		server(config)
	}

	if *config.RunClient {
		if *config.Destination == "-" {
			log.Fatalln("Client cannot write to stdout")
		}
		client(config)
	}
}

