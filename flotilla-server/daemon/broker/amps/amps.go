package amps

import (
	"fmt"
	"github.com/60East/amps-client-go"
	"github.com/tylertreat/Flotilla/flotilla-server/daemon/broker"
)

type Peer struct {
	client    *amps.Client
	readchan  chan []byte
	wrchan    chan []byte
	wrmsg     amps.Message
	errchan   chan error
	ackmsg    amps.Message
	ackcnt    int
	bookmarks string
	count     int
}

const topic = "queue"

func NewPeer(host string) (*Peer, error) {
	p := &Peer{}
	p.client = amps.NewClient("client")
	p.client.MessageHandler = p
	p.client.ErrorHandler = p
	err := p.client.Connect("localhost:5000")
	for err != nil {
		err = p.client.Connect("localhost:5000")
	}
	if err == nil {
		p.wrchan = make(chan []byte, 4096)
		p.readchan = make(chan []byte, 4096)
		p.errchan = make(chan error)
		p.wrmsg.Command = "p"
		p.wrmsg.Topic = "queue"
	} else {
		fmt.Printf("PEER ERROR %s\n", err)
	}
	go func() {
		for {
			data := <-p.wrchan
			p.wrmsg.Data = string(data)
			err := p.client.Send(&p.wrmsg)
			if err != nil {
				break
			}
		}
	}()

	// logon
	logonMessage := amps.Message{
		Command: "logon",
		AckType: "processed",
		UserId:  broker.GenerateName(),
	}
	p.client.Send(&logonMessage)
	<-p.readchan

	return p, nil
}
func (p *Peer) Subscribe() error {
	fmt.Println("SUBSCRIBE")
	m := &amps.Message{
		Command: "subscribe",
		Topic:   "queue",
		Options: "max_backlog=100",
		AckType: "processed",
	}
	p.ackmsg.Command = "sow_delete"
	p.ackmsg.Topic = "queue"
	p.client.Send(m)
	<-p.readchan
	return nil
}

func (p *Peer) Recv() ([]byte, error) {
	b := <-p.readchan
	return b, nil
}

func (p *Peer) Send() chan<- []byte {
	fmt.Println("SENDCHAN")
	return p.wrchan
}

func (p *Peer) Errors() <-chan error {
	return p.errchan
}

func (p *Peer) Done() {
}

func (p *Peer) Setup() {
	fmt.Println("SETUP")
}

func (p *Peer) Invoke(m *amps.Message) {
	p.count++
	data := make([]byte, len(m.Data))
	copy(data, m.Data)
	p.readchan <- data
	p.bookmarks += m.Bookmark
	p.ackcnt++
	if p.ackcnt > 80 {
		p.sendAcks()
	}
}

func (p *Peer) Error(e error) {
	p.errchan <- e
}

func (p *Peer) Teardown() {
}

func (p *Peer) sendAcks() {
	p.ackmsg.Bookmark = p.bookmarks
	p.client.Send(&p.ackmsg)
	p.bookmarks = ""
	p.ackcnt = 0
}

/*
	// Setup prepares the peer for testing.
	Setup()

	// Teardown performs any cleanup logic that needs to be performed after the
	// test is complete.
	Teardown()
*/
