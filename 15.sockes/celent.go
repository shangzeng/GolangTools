package main

import (
   "bufio"
   "io"
   "log"
   "net"
)


var (
   // 本地需要暴露的服务端口 
   localServerAddr = "127.0.0.1:8000"
   // 连接的IP端口
   remoteIP = "124.70.82.229"
   // 远端的服务控制通道，用来传递控制信息，如出现新连接和心跳
   remoteControlAddr = remoteIP + ":8009"
   // 远端服务端口，用来建立隧道
   remoteServerAddr  = remoteIP + ":8008"

   // 这里用于网络连接的 表示连接和存活显示
   KeepAlive     = "KEEP_ALIVE"
   NewConnection = "NEW_CONNECTION"
)


/*
客户端数据连接
*/
func main() {
   // 连接目标控制隧道信息 返回连接类型
   tcpConn, err := CreateTCPConn(remoteControlAddr)
   if err != nil {
      log.Println("[连接失败]" + remoteControlAddr + err.Error())
      return
   }
   log.Println("[已连接]" + remoteControlAddr)

   // 尝试读取连接数据 tcpConn
   reader := bufio.NewReader(tcpConn)
   for {
      // 以'\n'为结束符读入一行
      s, err := reader.ReadString('\n')
      if err != nil || err == io.EOF {
         break
      }

      // 当有新连接信号出现时，新建一个tcp连接
      if s == NewConnection+"\n" {
         // 将两个端口的信息互换
         go connectLocalAndRemote()
      }
   }

   log.Println("[已断开]" + remoteControlAddr)
}
// 新建一个TCP连接
func connectLocalAndRemote() {
   // 这里建立的是本地80端口的连接
   local := connectLocal()
   //
   remote := connectRemote()

   if local != nil && remote != nil {
      Join2Conn(local, remote)
   } else {
      if local != nil {
         _ = local.Close()
      }
      if remote != nil {
         _ = remote.Close()
      }
   }
}
// 127.0.0.1:8000 针对本地端口进行监听
func connectLocal() *net.TCPConn {
   conn, err := CreateTCPConn(localServerAddr)
   if err != nil {
      log.Println("[连接本地服务失败]" + err.Error())
   }
   return conn
}

// 127.0.0.1:8008 这里是要转发的端口来南方肌肉
func connectRemote() *net.TCPConn {
   conn, err := CreateTCPConn(remoteServerAddr)
   if err != nil {
      log.Println("[连接远端服务失败]" + err.Error())
   }
   return conn
}











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


// 测试目标IP是否开启   124.70.82.229:8009
// 如果开启返回一个连接
func CreateTCPConn(addr string) (*net.TCPConn, error) {
   // 一种写入的地址格式
   tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
   if err != nil {
      return nil, err
   }
   // 进行连接，返回 net.TCPConn 类型
   tcpListener, err := net.DialTCP("tcp",nil, tcpAddr)
   if err != nil {
      return nil, err
   }
   return tcpListener, nil
}

// 互换
func Join2Conn(local *net.TCPConn, remote *net.TCPConn) {
   go joinConn(local, remote)
   go joinConn(remote, local)
}
// 互换
func joinConn(local *net.TCPConn, remote *net.TCPConn) {
   defer local.Close()
   defer remote.Close()
   _, err := io.Copy(local, remote)
   if err != nil {
      log.Println("copy failed ", err.Error())
      return
   }
}