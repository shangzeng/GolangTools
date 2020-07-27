package functions


import(
  "os"
  "fmt"
  "strings"
  "strconv"
  "bufio"
  "encoding/json"
)





func ReadLines(path string) ([]string, int, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil,0, err
  }
  defer file.Close()

  var lines []string
  linecount :=0
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
    linecount++
  }
  return lines,linecount,scanner.Err()
}





func Configing(files string)(string, string ,string ,int ,int, string) {
  file, err := os.Open(files)
  if err != nil {
    fmt.Println("配置文件读取错误 ！")
    os.Exit(0)
  }
  defer file.Close()
  decoder := json.NewDecoder(file)
  conf := configuration{}
  err1 := decoder.Decode(&conf)
  if err1 != nil {
    fmt.Println("编码错误 ！请检查json文件格式 ！")
    os.Exit(0)
  }
  return conf.Ipfile,conf.Nmapports,conf.Nmapresults,conf.Nmapthreads,conf.Nmaptimeout, conf.Nmapresultfile
}




// 用于写入数据
func tracefile(str_content,name string)  {

    fd,_:=os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
    fd_content:=str_content
    buf:=[]byte(fd_content)
    fd.Write(buf)
    fd.Close()
}













// 读取文件数据，返回
func ReadIpList(fileName string) (ipList []Service) {
  ipListFile, err := os.Open(fileName)
  if err != nil {
    fmt.Println("Open ip List file err, %v", err)
    os.Exit(0)
  }

  defer ipListFile.Close()

  scanner := bufio.NewScanner(ipListFile)
  scanner.Split(bufio.ScanLines)

  for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
      continue
    }
    ipPort := strings.TrimSpace(line)
    t := strings.Split(ipPort, ":")
    t1 := strings.Split(ipPort, "\"")
    ip := t1[1]
    portProtocol := t[1]
    tmpPort := strings.Split(portProtocol, "|")
    port, _ := strconv.Atoi(tmpPort[0])
    protocol := strings.ToUpper(tmpPort[1])

    addr := Service{Ip: ip, Port: port, Protocol: protocol}
    ipList = append(ipList, addr)
  }

  return ipList
}




//  装载用户名
func ReadUserDict(userDict string) (users []string) {
  file, err := os.Open(userDict)
  if err != nil {
    fmt.Println("Open user dict file err", err)
    os.Exit(0)
  }

  defer file.Close()

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  for scanner.Scan() {
    user := strings.TrimSpace(scanner.Text())
    if user != "" {
      users = append(users, user)
    }
  }
  return users
}




// 装载密码
func ReadPasswordDict(passDict string) (password []string) {
  file, err := os.Open(passDict)
  if err != nil {
    fmt.Println("Open password dict file err", err)
    os.Exit(0)
  }

  defer file.Close()

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  for scanner.Scan() {
    passwd := strings.TrimSpace(scanner.Text())
    if passwd != "" {
      password = append(password, passwd)
    }
  }
  password = append(password, "")
  return password
}


// 将账号密码装载进入任务进行扫描
func AddUserAndPass (aliveIpList []Service, users []string, passwords []string) (aliveIpListtask []ServiceWeakPassword){
  for _, user := range users {
    for _, password := range passwords {
      for _, addr := range aliveIpList {
        service := ServiceWeakPassword{Server: addr, Username: user, Password: password}
        aliveIpListtask = append(aliveIpListtask, service)
      }
    }
  }
  return aliveIpListtask
}



