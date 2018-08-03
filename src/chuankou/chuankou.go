/*
 * This application can be used to experiment and test various serial port options
 */

package chuankou

import (
 "encoding/hex"
 "flag"
 "fmt"
 "io"
 "os"
 "jinzhi"
 "github.com/jacobsa/go-serial/serial"
 "github.com/larspensjo/config"
)

func usage() {
    fmt.Println("go-serial-test usage:")
    flag.PrintDefaults()
    os.Exit(-1)
}

func Read()(int,int) {
    c, _ := config.ReadDefault("config.ini")
    kahaoweishu,_:=c.Int("chuankou", "kahaoweishu")
    port1,_:=c.String("chuankou", "port")
    baud1,_:=c.Int("chuankou", "baud")
    txData1,_:=c.String("chuankou", "txData")
    stopbits1,_:=c.Int("chuankou", "stopbits")
    databits1,_:=c.Int("chuankou", "databits")

    port := flag.String("port", port1, "serial port to test (/dev/tty.wchusbserial1420, etc)")
    baud := flag.Uint("baud", uint(baud1), "Baud rate")
    txData := flag.String("txdata",txData1 , "data to send in hex format (01ab238b)")
    even := flag.Bool("even", false, "enable even parity")
    odd := flag.Bool("odd", false, "enable odd parity")
    rs485 := flag.Bool("rs485", false, "enable RS485 RTS for direction control")
    rs485HighDuringSend := flag.Bool("rs485_high_during_send", false, "RTS signal should be high during send")
    rs485HighAfterSend := flag.Bool("rs485_high_after_send", false, "RTS signal should be high after send")
    stopbits := flag.Uint("stopbits", uint(stopbits1), "Stop bits")
    databits := flag.Uint("databits", uint(databits1), "Data bits")
    chartimeout := flag.Uint("chartimeout", 100, "Inter Character timeout (ms)")
    minread := flag.Uint("minread", 0, "Minimum read count")
    rx := flag.Bool("rx", true, "Read data received")

    flag.Parse()

    if *port == "" {
        fmt.Println("Must specify port")
        usage()
    }

    if *even && *odd {
        fmt.Println("can't specify both even and odd parity")
        usage()
    }

    parity := serial.PARITY_NONE

    if *even {
        parity = serial.PARITY_EVEN
    } else if *odd {
        parity = serial.PARITY_ODD
    }

    options := serial.OpenOptions{
        PortName:               *port,
        BaudRate:               *baud,
        DataBits:               *databits,
        StopBits:               *stopbits,
        MinimumReadSize:        *minread,
        InterCharacterTimeout:  *chartimeout,
        ParityMode:             parity,
        Rs485Enable:            *rs485,
        Rs485RtsHighDuringSend: *rs485HighDuringSend,
        Rs485RtsHighAfterSend:  *rs485HighAfterSend,
    }

    f, err := serial.Open(options)

    if err != nil {
        fmt.Println("Error opening serial port: ", err)
        os.Exit(-1)
    } else {
        defer f.Close()
    }

    if *txData != "" {
        txData_, err := hex.DecodeString(*txData)

        if err != nil {
            fmt.Println("Error decoding hex data: ", err)
            os.Exit(-1)
        }

        // fmt.Println("Sending: ", hex.EncodeToString(txData_))

        f.Write(txData_)

        // if err != nil {
        //     fmt.Println("Error writing to serial port: ", err)
        // } else {
        //     // fmt.Printf("Wrote %v bytes\n", count)
        // }

    }
    var bufComplete string
    if *rx {
        for {
            buf := make([]byte, 32)
            n, err := f.Read(buf)

            if err != nil {
                if err != io.EOF {
                    fmt.Println("Error reading from serial port: ", err)
                }
            } else {
                buf = buf[:n]
                if len(buf)>0 {
                    bufComplete += hex.EncodeToString(buf)

                    if len(bufComplete)>=kahaoweishu {
                        sjzdyg:= string([]byte(bufComplete)[8:])
                        sjzdyg = string([]byte(sjzdyg)[:2])
                        sjzdeg:= string([]byte(bufComplete)[10:])
                        sjzdeg = string([]byte(sjzdeg)[:4])
                        // fmt.Println(sjzdyg)//aa080400 d4 2804 f4bb
                        // fmt.Println(sjzdeg)//aa080400 d4 2804 f4bb


                        // fmt.Println(jinzhi.AnyToDecimal(string([]byte(bufComplete)[8:]),16))//aa080400 d4 2804 f4bb
                        f.Close()
                        return jinzhi.AnyToDecimal(sjzdyg,16),jinzhi.AnyToDecimal(sjzdeg,16)
                        break

                    }

                    
                }else{
                    if len(bufComplete)>=kahaoweishu {
                        sjzdyg:= string([]byte(bufComplete)[8:])
                        sjzdyg = string([]byte(sjzdyg)[:2])
                        sjzdeg:= string([]byte(bufComplete)[10:])
                        sjzdeg = string([]byte(sjzdeg)[:4])
                        // fmt.Println(jinzhi.AnyToDecimal(bufComplete,16))//aa080400 d4 2804 f4bb
                        f.Close()
                        return jinzhi.AnyToDecimal(sjzdyg,16),jinzhi.AnyToDecimal(sjzdeg,16)
                        break
                    }
                }
                
                
            }
        }
    }
    return -1,-1
}
