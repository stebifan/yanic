package yanic

/**
 * This database type is for injecting into another yanic instance.
 */
import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	config Config
	conn   net.Conn
}

type Config map[string]interface{}

func (c Config) Address() string {
	return c["address"].(string)
}

func init() {
	database.RegisterAdapter("yanic", Connect)
}

func Connect(configuration map[string]interface{}) (database.Connection, error) {
	var config Config
	config = configuration

	conn, err := net.Dial("udp6", config.Address())
	if err != nil {
		return nil, err
	}

	return &Connection{conn: conn, config: config}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	res := &data.ResponseData{
		NodeInfo:   node.Nodeinfo,
		Statistics: node.Statistics,
		Neighbours: node.Neighbours,
	}
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	flater, err := flate.NewWriter(writer, flate.NoCompression)
	if err != nil {
		log.Printf("[database-yanic] could not create flater: %s", err)
		return
	}
	defer flater.Close()
	err = json.NewEncoder(flater).Encode(res)
	if err != nil {
		nodeid := "unknown"
		if node.Nodeinfo != nil && node.Nodeinfo.NodeID != "" {
			nodeid = node.Nodeinfo.NodeID
		}
		log.Printf("[database-yanic] could not send %s node: %s", nodeid, err)
		return
	}
	flater.Flush()
	writer.Flush()
	conn.conn.Write(b.Bytes())

}

func (conn *Connection) InsertLink(link *runtime.Link, time time.Time) {
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string) {
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
}

func (conn *Connection) Close() {
	conn.conn.Close()
}
