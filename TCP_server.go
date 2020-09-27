package main

import (
  "net"
  "math/rand"
  "time"
  "bufio"
  b64 "encoding/base64"
  "log"
)

type server struct{            
  clients []net.Conn
  key []byte
  listener net.Listener
}

func (server *server) handle_clients(sender_conn net.Conn) {
  for {
    reader := bufio.NewReader(sender_conn)

    message, err := reader.ReadString('\n')
    error_handler(err)
    // gets message
  
    server.broadcast(sender_conn, message)
    // sends message
  }
}

func (server *server) broadcast(sender net.Conn, message string) {
  for _, client := range server.clients {
    if client != sender {
      _, err := client.Write([]byte(message))

      if err != nil {
        log.Debug(err)
        continue
      }
    }
  }
}

func (server *server) get_new_clients() net.Conn {

  new_client, err := server.listener.Accept()
  error_handler(err)

  if new_client != nil {
    server.clients = append(server.clients, new_client)
    // adds clien to list

    _, err2 := new_client.Write(server.key)
    error_handler(err2)
    // sends encryption key
  
    log.Println("Serving: ", new_client.RemoteAddr())
  }
  return new_client
}

func (server *server) shutdown() {
  for _, client := range server.clients {
    client.Close()
  }
}

func main() {

  key := gen_byte_string(32)
  key = base64(key)

  tcpAddr, err := net.ResolveTCPAddr("tcp", ":6545")
  error_handler(err)

  listening_socket, err := net.ListenTCP("tcp", tcpAddr)

  server := server{
    clients: make([]net.Conn, 0),
    key: key,
    listener: listening_socket,
  }

  for {
    client := server.get_new_clients()
    if client != nil {
      go server.handle_clients(client)
    }
  }
  server.shutdown()
}

func base64(message []byte) []byte {
	b := make([]byte, b64.StdEncoding.EncodedLen(len(message)))
	b64.StdEncoding.Encode(b, message)
	return b
}

func gen_byte_string(length int) []byte {
  b := make([]byte, length)

  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return b
}

func error_handler(err error) {
  if err != nil {
    log.Debug(err)
  }
}
