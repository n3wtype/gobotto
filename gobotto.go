package main

import "os"
import "fmt"
import "net"
import "bufio"

type ircBot struct {
	conn net.Conn
	nick string
	channel string
}

func NewIrcBot (server, port, channel, nick string) *ircBot {

	conn, err := net.Dial ("tcp", fmt.Sprintf("%s:%s", server, port))
	if err==nil{}
	fmt.Fprintf (conn, fmt.Sprintf ("NICK %s\r\n", nick))
	fmt.Fprintf (conn, fmt.Sprintf ("USER %s 8 *: %s\r\n", nick, nick))
	fmt.Fprintf (conn, fmt.Sprintf ("JOIN %s\r\n", channel))
	
	return &ircBot{nick: nick, conn: conn, channel: channel}
}

func (i ircBot) Say(msg string) {
	fmt.Fprintf (i.conn, fmt.Sprintf ("PRIVMSG %s :%s\r\n", i.channel, msg))
}

func (i ircBot) ReadLine () string {
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

