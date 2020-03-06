// Autogenerated by Thrift Compiler (0.13.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"github.com/apache/thrift/lib/go/thrift"
	"openstars/core/bigset/listi64"
)

var _ = listi64.GoUnusedProtection__

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  TPutItemResult bsgPutItem(TContainerKey rootID, TItem item)")
  fmt.Fprintln(os.Stderr, "  bool bsgRemoveItem(TMetaKey key, TItemKey itemKey)")
  fmt.Fprintln(os.Stderr, "  TExistedResult bsgExisted(TContainerKey rootID, TItemKey itemKey)")
  fmt.Fprintln(os.Stderr, "  TItemResult bsgGetItem(TContainerKey rootID, TItemKey itemKey)")
  fmt.Fprintln(os.Stderr, "  TItemSetResult bsgGetSlice(TContainerKey rootID, i32 fromIDX, i32 count)")
  fmt.Fprintln(os.Stderr, "  TItemSetResult bsgGetSliceFromItem(TContainerKey rootID, TItemKey fromKey, i32 count)")
  fmt.Fprintln(os.Stderr, "  TItemSetResult bsgGetSliceR(TContainerKey rootID, i32 fromIDX, i32 count)")
  fmt.Fprintln(os.Stderr, "  TItemSetResult bsgGetSliceFromItemR(TContainerKey rootID, TItemKey fromKey, i32 count)")
  fmt.Fprintln(os.Stderr, "  TSplitBigSetResult splitBigSet(TContainerKey rootID, TContainerKey brotherRootID, i64 currentSize)")
  fmt.Fprintln(os.Stderr, "  TItemSetResult bsgRangeQuery(TContainerKey rootID, TItemKey startKey, TItemKey endKey)")
  fmt.Fprintln(os.Stderr, "  bool bsgBulkLoad(TContainerKey rootID, TItemSet setData)")
  fmt.Fprintln(os.Stderr, "  TMultiPutItemResult bsgMultiPut(TContainerKey rootID, TItemSet setData, bool getAddedItems, bool getReplacedItems)")
  fmt.Fprintln(os.Stderr, "  TBigSetLI64Data getSetGenData(TMetaKey metaID)")
  fmt.Fprintln(os.Stderr, "  void putSetGenData(TMetaKey metaID, TBigSetLI64Data metadata)")
  fmt.Fprintln(os.Stderr, "  i64 getTotalCount(TMetaKey metaID)")
  fmt.Fprintln(os.Stderr, "  i64 removeAll(TContainerKey rootID)")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

type httpHeaders map[string]string

func (h httpHeaders) String() string {
  var m map[string]string = h
  return fmt.Sprintf("%s", m)
}

func (h httpHeaders) Set(value string) error {
  parts := strings.Split(value, ": ")
  if len(parts) != 2 {
    return fmt.Errorf("header should be of format 'Key: Value'")
  }
  h[parts[0]] = parts[1]
  return nil
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  headers := make(httpHeaders)
  var parsedUrl *url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Var(headers, "H", "Headers to set on the http(s) request (e.g. -H \"Key: Value\")")
  flag.Parse()
  
  if len(urlString) > 0 {
    var err error
    parsedUrl, err = url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
    if len(headers) > 0 {
      httptrans := trans.(*thrift.THttpClient)
      for key, value := range headers {
        httptrans.SetHeader(key, value)
      }
    }
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  iprot := protocolFactory.GetProtocol(trans)
  oprot := protocolFactory.GetProtocol(trans)
  client := listi64.NewTBSLI64DataServiceClient(thrift.NewTStandardClient(iprot, oprot))
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "bsgPutItem":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "BsgPutItem requires 2 args")
      flag.Usage()
    }
    argvalue0, err51 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err51 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    arg52 := flag.Arg(2)
    mbTrans53 := thrift.NewTMemoryBufferLen(len(arg52))
    defer mbTrans53.Close()
    _, err54 := mbTrans53.WriteString(arg52)
    if err54 != nil {
      Usage()
      return
    }
    factory55 := thrift.NewTJSONProtocolFactory()
    jsProt56 := factory55.GetProtocol(mbTrans53)
    argvalue1 := listi64.NewTItem()
    err57 := argvalue1.Read(jsProt56)
    if err57 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.BsgPutItem(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "bsgRemoveItem":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "BsgRemoveItem requires 2 args")
      flag.Usage()
    }
    argvalue0, err58 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err58 != nil {
      Usage()
      return
    }
    value0 := listi64.TMetaKey(argvalue0)
    argvalue1 := []byte(flag.Arg(2))
    value1 := listi64.TItemKey(argvalue1)
    fmt.Print(client.BsgRemoveItem(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "bsgExisted":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "BsgExisted requires 2 args")
      flag.Usage()
    }
    argvalue0, err60 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err60 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    argvalue1 := []byte(flag.Arg(2))
    value1 := listi64.TItemKey(argvalue1)
    fmt.Print(client.BsgExisted(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "bsgGetItem":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "BsgGetItem requires 2 args")
      flag.Usage()
    }
    argvalue0, err62 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err62 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    argvalue1 := []byte(flag.Arg(2))
    value1 := listi64.TItemKey(argvalue1)
    fmt.Print(client.BsgGetItem(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "bsgGetSlice":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "BsgGetSlice requires 3 args")
      flag.Usage()
    }
    argvalue0, err64 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err64 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    tmp1, err65 := (strconv.Atoi(flag.Arg(2)))
    if err65 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    tmp2, err66 := (strconv.Atoi(flag.Arg(3)))
    if err66 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.BsgGetSlice(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "bsgGetSliceFromItem":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "BsgGetSliceFromItem requires 3 args")
      flag.Usage()
    }
    argvalue0, err67 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err67 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    argvalue1 := []byte(flag.Arg(2))
    value1 := listi64.TItemKey(argvalue1)
    tmp2, err69 := (strconv.Atoi(flag.Arg(3)))
    if err69 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.BsgGetSliceFromItem(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "bsgGetSliceR":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "BsgGetSliceR requires 3 args")
      flag.Usage()
    }
    argvalue0, err70 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err70 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    tmp1, err71 := (strconv.Atoi(flag.Arg(2)))
    if err71 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    tmp2, err72 := (strconv.Atoi(flag.Arg(3)))
    if err72 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.BsgGetSliceR(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "bsgGetSliceFromItemR":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "BsgGetSliceFromItemR requires 3 args")
      flag.Usage()
    }
    argvalue0, err73 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err73 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    argvalue1 := []byte(flag.Arg(2))
    value1 := listi64.TItemKey(argvalue1)
    tmp2, err75 := (strconv.Atoi(flag.Arg(3)))
    if err75 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.BsgGetSliceFromItemR(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "splitBigSet":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "SplitBigSet requires 3 args")
      flag.Usage()
    }
    argvalue0, err76 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err76 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    argvalue1, err77 := (strconv.ParseInt(flag.Arg(2), 10, 64))
    if err77 != nil {
      Usage()
      return
    }
    value1 := listi64.TContainerKey(argvalue1)
    argvalue2, err78 := (strconv.ParseInt(flag.Arg(3), 10, 64))
    if err78 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    fmt.Print(client.SplitBigSet(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "bsgRangeQuery":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "BsgRangeQuery requires 3 args")
      flag.Usage()
    }
    argvalue0, err79 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err79 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    argvalue1 := []byte(flag.Arg(2))
    value1 := listi64.TItemKey(argvalue1)
    argvalue2 := []byte(flag.Arg(3))
    value2 := listi64.TItemKey(argvalue2)
    fmt.Print(client.BsgRangeQuery(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "bsgBulkLoad":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "BsgBulkLoad requires 2 args")
      flag.Usage()
    }
    argvalue0, err82 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err82 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    arg83 := flag.Arg(2)
    mbTrans84 := thrift.NewTMemoryBufferLen(len(arg83))
    defer mbTrans84.Close()
    _, err85 := mbTrans84.WriteString(arg83)
    if err85 != nil {
      Usage()
      return
    }
    factory86 := thrift.NewTJSONProtocolFactory()
    jsProt87 := factory86.GetProtocol(mbTrans84)
    argvalue1 := listi64.NewTItemSet()
    err88 := argvalue1.Read(jsProt87)
    if err88 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.BsgBulkLoad(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "bsgMultiPut":
    if flag.NArg() - 1 != 4 {
      fmt.Fprintln(os.Stderr, "BsgMultiPut requires 4 args")
      flag.Usage()
    }
    argvalue0, err89 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err89 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    arg90 := flag.Arg(2)
    mbTrans91 := thrift.NewTMemoryBufferLen(len(arg90))
    defer mbTrans91.Close()
    _, err92 := mbTrans91.WriteString(arg90)
    if err92 != nil {
      Usage()
      return
    }
    factory93 := thrift.NewTJSONProtocolFactory()
    jsProt94 := factory93.GetProtocol(mbTrans91)
    argvalue1 := listi64.NewTItemSet()
    err95 := argvalue1.Read(jsProt94)
    if err95 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    argvalue2 := flag.Arg(3) == "true"
    value2 := argvalue2
    argvalue3 := flag.Arg(4) == "true"
    value3 := argvalue3
    fmt.Print(client.BsgMultiPut(context.Background(), value0, value1, value2, value3))
    fmt.Print("\n")
    break
  case "getSetGenData":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetSetGenData requires 1 args")
      flag.Usage()
    }
    argvalue0, err98 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err98 != nil {
      Usage()
      return
    }
    value0 := listi64.TMetaKey(argvalue0)
    fmt.Print(client.GetSetGenData(context.Background(), value0))
    fmt.Print("\n")
    break
  case "putSetGenData":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "PutSetGenData requires 2 args")
      flag.Usage()
    }
    argvalue0, err99 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err99 != nil {
      Usage()
      return
    }
    value0 := listi64.TMetaKey(argvalue0)
    arg100 := flag.Arg(2)
    mbTrans101 := thrift.NewTMemoryBufferLen(len(arg100))
    defer mbTrans101.Close()
    _, err102 := mbTrans101.WriteString(arg100)
    if err102 != nil {
      Usage()
      return
    }
    factory103 := thrift.NewTJSONProtocolFactory()
    jsProt104 := factory103.GetProtocol(mbTrans101)
    argvalue1 := listi64.NewTBigSetLI64Data()
    err105 := argvalue1.Read(jsProt104)
    if err105 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.PutSetGenData(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "getTotalCount":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetTotalCount requires 1 args")
      flag.Usage()
    }
    argvalue0, err106 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err106 != nil {
      Usage()
      return
    }
    value0 := listi64.TMetaKey(argvalue0)
    fmt.Print(client.GetTotalCount(context.Background(), value0))
    fmt.Print("\n")
    break
  case "removeAll":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "RemoveAll requires 1 args")
      flag.Usage()
    }
    argvalue0, err107 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err107 != nil {
      Usage()
      return
    }
    value0 := listi64.TContainerKey(argvalue0)
    fmt.Print(client.RemoveAll(context.Background(), value0))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}