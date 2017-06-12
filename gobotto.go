package main

import "os"
import "fmt"
import "net"
import "bufio"
import "bytes"
import "./irc"

const msg_buf_size = 10240

type ircBot struct {
	conn net.Conn
	nick string
	channel string
	conn_params connParams
}

type ircMsg struct {
	prefix string
	command string
	params []string
}

type ircMsgBroker struct {

}

type connParams struct {
	port string
	server string
}

func parseMessage (line string) *ircMsg {

	return nil
}

func NewIrcBot (server, port, channel, nick string) *ircBot {

	cp := connParams {port: port, server: server}
	i := ircBot {nick: nick, conn: nil, channel: channel, conn_params: cp }

	return &i
}

func (self ircBot) Say(msg string) {
	fmt.Fprintf (self.conn, fmt.Sprintf ("PRIVMSG %s :%s\r\n", self.channel, msg))
}

func (i *ircBot) reconnect() {

	conn, err := net.Dial ("tcp", fmt.Sprintf("%s:%s", i.conn_params.server, i.conn_params.port))

	if err == nil {}
	i.conn = conn

}

func (i ircBot) ircSetNick (nick string) {

	i.nick = nick

	fmt.Fprintf (i.conn, fmt.Sprintf ("NICK %s\r\n", i.nick))
	fmt.Fprintf (i.conn, fmt.Sprintf ("USER %s 8 *: %s\r\n", i.nick, i.nick))
}

func (i ircBot) ircJoinChannel (channel string) {
	fmt.Fprintf (i.conn, fmt.Sprintf ("JOIN %s\r\n", channel))
}

func splitCRLF (data []byte, atEOF bool) (advance int, token []byte, err error){ 
	if atEOF && len (data) == 0{
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte{'\r','\n'}); i >= 0 {
		return i + 2, data[0:i], nil
	}

	if atEOF {
            return len(data), data, nil
    }

	return 0, nil, nil
}

func (i *ircBot) ReadLine (done chan bool) {
	rd := bufio.NewReaderSize(i.conn, msg_buf_size)
	sc := bufio.NewScanner(rd)
	sc.Split(splitCRLF)

	for sc.Scan() {
			fmt.Println(sc.Text())
	}

}


func (i *ircBot) run (done chan bool) {
	i.reconnect()
	i.ircSetNick(i.nick)
	i.ircJoinChannel(i.channel)
	i.ReadLine(done)
}

func main() {
	nick := GetParamValue("--nick")
	server := GetParamValue("--server")
	channel := GetParamValue("--channel")
	port := GetParamValue("--port")

	fmt.Printf (fmt.Sprintf ("DEBUG (nick: %s. server: %s, channel: %s, port: %s\r\n", 
		nick, server, channel, port ))
	fmt.Println ("")

	bot := NewIrcBot (server, port, channel, nick)

	done := make (chan bool,1)
	go bot.run(done)

	<-done

 }

func SliceIndex (limit int, predicate func (i int) bool ) int {
	for i := 0; i < limit; i++ {
        if predicate(i) {
            return i
        }
    }
    return -1
}

func FindParam (param string) int {
	return SliceIndex (len(os.Args), func (i int)  bool { return os.Args[i] == param })
}

func GetParamValue (param string) string{
	if k := FindParam(param); k > 0 {
		if v := k+1; v <= len(os.Args)  { 
			return os.Args[v]
		} else {
			return ""
		}
	} else {
		return ""
	}
}

