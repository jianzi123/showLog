package util

import (
	"io"
	"log"
	"os/exec"
	"fmt"
	"bufio"
	"golang.org/x/net/websocket"
	"sync"
)
type SockCon struct {
        Bflag bool
        User string
        Job string
		Content string
        Id int  
	MapStr map[string]*websocket.Conn
        Mut * sync.RWMutex
}

func NewSock() (*SockCon, error) {
	s := &SockCon{
		Bflag: 		false,
		User:		"",
		Job:		"",
		Content:	"",
		Id:			0,
		MapStr: 	make(map[string]*websocket.Conn),
		Mut:	 	new(sync.RWMutex),
	}	
	return s, nil
}

func (s *SockCon)ExecCmd(cmdName string, params []string) {
	cmd := exec.Command(cmdName, params...)
	fmt.Println(cmd.Args)
	var logMsg string
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)	
	}	
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}	
	log.Print("waiting for command to finish.")
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		WriteFile("./", "jion.log", line)
		s.Mut.Lock()
		if v, ok := s.MapStr["jion"]; ok {
			if err = websocket.Message.Send(v, line); err != nil {
                        	fmt.Println(err)
            }
		}else{
			logMsg = logMsg + line
			s.Content = logMsg
		} 
		s.Mut.Unlock()
		fmt.Println(line)
	}
	err1 := cmd.Wait()
	if err1 != nil {
		log.Print( err1)
		return
	}
}

