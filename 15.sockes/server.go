package main

import (
   "io"
   "log"
   "net"
   "strconv"
   "sync"
   "time"

)

const (
   controlAddr = "0.0.0.0:8009"
   tunnelAddr  = "0.0.0.0:8008"
   visitAddr   = "0.0.0.0:8007"
   KeepAlive     = "KEEP_ALIVE"
   NewConnection = "NEW_CONNECTION"
)

var (
   clientConn         *net.TCPConn
   connectionPool     map[string]*ConnMatch
   connectionPoolLock sync.Mutex
)

type ConnMatch struct {
   addTime time.Time
   accept  *net.TCPConn
}


/*
服务端数据连接
*/
func main() {
   connectionPool = make(map[string]*ConnMatch, 32)
   go createControlChannel()                                  // 8009 端口，监听链接，防止断开
   go acceptUserRequest()                                     // 0.0.0.0:8007  用于发送消息（这里貌似没怎用）
   go acceptClientRequest()                                   // 0.0.0.0:8008  将获取的数据进行交换
   cleanConnectionPool()                                      // 关于链接池
}

// 创建一个控制通道，用于传递控制消息，如：心跳，创建新连接
func createControlChannel() {
   // 8009 端口获取链接请求
   tcpListener, err := CreateTCPListener(controlAddr)
   if err != nil {
      panic(err)
   }

   log.Println("[已监听]" + controlAddr)
   for {
      tcpConn, err := tcpListener.AcceptTCP()
      if err != nil {
         log.Println(err)
         continue
      }

      log.Println("[新连接]" + tcpConn.RemoteAddr().String())
      // 如果当前已经有一个客户端存在，则丢弃这个链接
      if clientConn != nil {
         _ = tcpConn.Close()
      } else {
         clientConn = tcpConn
         // 保持链接防止断掉
         go keepAlive()
      }
   }
}

// 和客户端保持一个心跳链接
func keepAlive() {
   go func() {
      for {
         if clientConn == nil {
            return
         }
         _, err := clientConn.Write(([]byte)(KeepAlive + "\n"))
         if err != nil {
            log.Println("[已断开客户端连接]", clientConn.RemoteAddr())
            clientConn = nil
            return
         }
         time.Sleep(time.Second * 3)
      }
   }()
}

// 监听来自用户的请求
func acceptUserRequest() {
   // 8007 
   tcpListener, err := CreateTCPListener(visitAddr)
   if err != nil {
      panic(err)
   }
   defer tcpListener.Close()
   for {
      tcpConn, err := tcpListener.AcceptTCP()
      if err != nil {
         continue
      }
      // 放入链接池 进行记录
      addConn2Pool(tcpConn)
      // 发送消息
      sendMessage(NewConnection + "\n")
   }
}

// 将用户来的连接放入连接池中
func addConn2Pool(accept *net.TCPConn) {
   connectionPoolLock.Lock()
   defer connectionPoolLock.Unlock()

   now := time.Now()
   connectionPool[strconv.FormatInt(now.UnixNano(), 10)] = &ConnMatch{now, accept,}
}

// 发送给客户端新消息
func sendMessage(message string) {
   if clientConn == nil {
      log.Println("[无已连接的客户端]")
      return
   }
   _, err := clientConn.Write([]byte(message))
   if err != nil {
      log.Println("[发送消息异常]: message: ", message)
   }
}

// 接收客户端来的请求并建立隧道
func acceptClientRequest() {
   tcpListener, err := CreateTCPListener(tunnelAddr)
   if err != nil {
      panic(err)
   }
   defer tcpListener.Close()

   for {
      tcpConn, err := tcpListener.AcceptTCP()
      if err != nil {
         continue
      }
      // 将数据进行交换
      go establishTunnel(tcpConn)
   }
}
// 将获取的数据进行交换
func establishTunnel(tunnel *net.TCPConn) {
   connectionPoolLock.Lock()
   defer connectionPoolLock.Unlock()

   for key, connMatch := range connectionPool {
      if connMatch.accept != nil {
         go Join2Conn(connMatch.accept, tunnel)
         delete(connectionPool, key)
         return
      }
   }

   _ = tunnel.Close()
}

func cleanConnectionPool() {
   for {
      connectionPoolLock.Lock()
      for key, connMatch := range connectionPool {
         if time.Now().Sub(connMatch.addTime) > time.Second*10 {
            _ = connMatch.accept.Close()
            delete(connectionPool, key)
         }
      }
      connectionPoolLock.Unlock()
      time.Sleep(5 * time.Second)
   }
}






// 创建监听，返回连接数据
func CreateTCPListener(addr string) (*net.TCPListener, error) {
   tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
   if err != nil {
       return nil, err
   }
   tcpListener, err := net.ListenTCP("tcp", tcpAddr)
   if err != nil {
       return nil, err
   }
   return tcpListener, nil
}

func CreateTCPConn(addr string) (*net.TCPConn, error) {
   tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
   if err != nil {
      return nil, err
   }
   tcpListener, err := net.DialTCP("tcp",nil, tcpAddr)
   if err != nil {
      return nil, err
   }
   return tcpListener, nil
}

func Join2Conn(local *net.TCPConn, remote *net.TCPConn) {
   go joinConn(local, remote)
   go joinConn(remote, local)
}

func joinConn(local *net.TCPConn, remote *net.TCPConn) {
   defer local.Close()
   defer remote.Close()
   _, err := io.Copy(local, remote)
   if err != nil {
      log.Println("copy failed ", err.Error())
      return
   }
}