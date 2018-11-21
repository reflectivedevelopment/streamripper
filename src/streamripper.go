package main

import "crypto/rand"
import "crypto/tls"
import "crypto/x509"
import "flag"
//import "fmt"
import "io/ioutil"
import "log"
//import "net"
import "time"

import "ripper"
/*
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		//fmt.Println("worker", id, "started  job", j)
		//time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
*/

func server(certPath *string, keyPath *string, hostName *string, port *string) {
	certBytes, err := ioutil.ReadFile(*certPath)

	if err != nil {
		log.Fatalln("Unable to read cert", err)
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Unable to add certificate to certificate pool")
	}

	cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)

	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Ensure that we only use our "CA" to validate certificates
		ClientCAs: clientCertPool,
		// TLS 1.2 because we can
		MinVersion: tls.VersionTLS12,
	}

	config.Rand = rand.Reader
	config.BuildNameToCertificate()

	ln, err := tls.Listen("tcp", *hostName + ":" + *port, &config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//go handleConnection(conn)
		res, err := conn.Write([]byte("bla bla"));
		if err != nil {
			log.Println(res, err)
			continue
		}
		time.Sleep(time.Second)
		log.Printf("Connection from %v closed.", conn.RemoteAddr())
		conn.Close()
	}
}

func client(certPath *string, keyPath *string, hostName *string, port *string) {
	certBytes, err := ioutil.ReadFile(*certPath)

	if err != nil {
		log.Fatalln("Unable to read cert", err)
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Unable to add certificate to certificate pool")
	}

	cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: clientCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	conn, err := tls.Dial("tcp", *hostName+":"+*port, tlsConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	data := make([]byte, 1024)
	conn.Read(data)
	log.Println(string(data[:]))
}

func main() {
/*	jobs := make(chan int, 100)
	results := make(chan int, 100)
	count := 200

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= count; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= count; a++ {
		r := <-results
		fmt.Println(r);
	}*/

	runServer := flag.Bool("server", false, "Run as Server")
	runClient := flag.Bool("client", false, "Run as Client")

	keyPath := flag.String("key", "mykey.pem", "Key")
	certPath := flag.String("cert", "mycert.pem", "Cert")

	hostName := flag.String("host", "localhost", "Hostname")
	port := flag.String("port", "8888", "Port")

	flag.Parse()

	if *runServer {
		server(certPath, keyPath, hostName, port)
	}

	if *runClient {
		client(certPath, keyPath, hostName, port)
	}
}

