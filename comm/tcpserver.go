package comm

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"runtime"
	"time"
	_ "strings" //..
	"bytes"
	"encoding/binary"

)

//Client ...
type Client struct {
   devid string
   conn *net.Conn
}

//TCPServer ...
type TCPServer struct {
	clients []*Client
}


func (p *TCPServer) remove(slice []*Client, elems *Client) []*Client {
    for i := range slice {
        if slice[i] == elems {
            slice = append(slice[:i], slice[i+1:]...) //打散传入 [ *Client[i+1] , *Client[i+2] ...]
            return  slice
        }
    }
    return  slice
}

// IntToBytes  this  is a 
func (p *TCPServer)IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
// BytesToInt is 
func (p *TCPServer)BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

//PrintClients is 
func (p *TCPServer) PrintClients() {

		fmt.Println("================================================");
		for _, v := range p.clients {	       
			conn1 :=*v.conn     
			ipStr := conn1.RemoteAddr().String()	
		    fmt.Println(v.devid,ipStr)
		}
		fmt.Println("================================================");

}
//Init is a 
func (p *TCPServer) Init(ipport string) {

	l, err := net.Listen("tcp",ipport )
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + ipport)

	for {
	// 接收一个client
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go p.handleRequest(conn)
	}
}

//Request5Seconds is 
func (p* TCPServer) Request5Seconds() {
	t1 := time.NewTimer(time.Second * 5)
	for {
		select {
        case <-t1.C:
        	go func(){
	            for _, v := range p.clients {	            	
	            		writer := bufio.NewWriter(*v.conn)
				        writer.Write([]byte("@T\r\n"))		        
				        writer.Flush()	            			
			    }
		    }()		
            t1.Reset(time.Second * 5)
        }
        runtime.Gosched()
    }
}

//Close is 
func (p* TCPServer)Close() {
	fmt.Println("Closed!")
}
func (p* TCPServer) handleRequest(conn net.Conn) {


	//conn.SetReadDeadline(time.Now().Add(time.Minute * 3));
    defer conn.Close();

	ipStr := conn.RemoteAddr().String()
	// 构建reader和writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	cli := &Client{"", &conn}
	p.clients = append(p.clients,cli);

	defer func() {
		fmt.Println("Disconnected :" + ipStr)		
		conn.Close()
		p.clients = p.remove(p.clients , cli)		
	}()

	writer.Write([]byte("@T\r\n"))		        
	writer.Flush()	   

	for {
		// 读取一行数据, 以"\n"结尾
		b, _, err := reader.ReadLine()
		if err != nil {
			return
		}		
		//encodedStr := hex.EncodeToString(b)
		str := string(b)
		fmt.Println("Process:",str)
		if str[0:1] == "!"{			
			cli.devid = str[1:4]
		}
	}

	
}


