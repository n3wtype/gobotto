package main

import "os"
import "fmt"
import "net"
import "bufio"

type ircBot struct {
	conn net.Conn
	nick string
	channel string
	conn_params connParams
}

type connParams struct {
	port string
	server string
}

func NewIrcBot (server, port, channel, nick string) *ircBot {

	cp := connParams {port: port, server: server}
	i := ircBot {nick: nick, conn: nil, channel: channel, conn_params: cp }

	i.reconnect()
	i.ircSetNick(nick)
	i.ircJoinChannel(channel)

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

func (i *ircBot) ReadLine () string {
	line, err := bufio.NewReader(i.conn).ReadString('\n')
	if err==nil {}
	return line
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
	//bot.Say("test")
	for true {
		fmt.Print(bot.ReadLine())
	}

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

