package marionette_client

import (
    "net"
    "fmt"
    "encoding/json"
    "strconv"
)


func main() {
    var conn, err = net.Dial("tcp", "127.0.0.1:2828")
    if err != nil {
        fmt.Println(err)
    }

    readBuf := make([]byte, 2048)
    var inter, _ = conn.Read(readBuf)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(inter);
    fmt.Println(string(readBuf));

    writeBuf := make(map[string]interface{})

    writeBuf["name"] = "newSession"
    writeBuf["parameters"] = map[string]interface{} { "sessionId": "111", "capabilities": nil}

    var mashalled, _ = json.Marshal(writeBuf)
    var endString = fmt.Sprintf("%s:%s", strconv.Itoa(len(string(mashalled))), string(mashalled))

    conn.Write([]byte(endString))
    fmt.Println(endString)

    inter, _ = conn.Read(readBuf)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(string(readBuf));
    conn.Close()
}
