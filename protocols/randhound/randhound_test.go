package randhound_test

import (
	"log"
	"testing"
	"time"

	"github.com/dedis/cothority/lib/dbg"
	"github.com/dedis/cothority/lib/sda"
	"github.com/dedis/cothority/protocols/randhound"
)

func TestRandHound(t *testing.T) {

	// Setup parameters
	var name string = "RandHound"             // Protocol name
	var nodes int = 10                        // Number of nodes (peers + leader)
	var trustees int = 5                      // Number of trustees
	var purpose string = "RandHound test run" // Purpose

	local := sda.NewLocalTest()
	_, _, tree := local.GenTree(nodes, false, true, true)
	defer local.CloseAll()

	dbg.TestOutput(testing.Verbose(), 1)

	// Setup and Start RandHound
	log.Printf("RandHound - starting")
	leader, err := local.CreateNewNodeName(name, tree)
	if err != nil {
		t.Fatal("Couldn't initialise RandHound protocol:", err)
	}
	rh := leader.ProtocolInstance().(*randhound.RandHound)
	err = rh.Setup(nodes, trustees, purpose)
	if err != nil {
		t.Fatal("Couldn't initialise RandHound protocol:", err)
	}
	log.Printf("RandHound - group config: %d %d %d %d %d %d\n", rh.Group.N, rh.Group.F, rh.Group.L, rh.Group.K, rh.Group.R, rh.Group.T)
	leader.StartProtocol()

	bytes := make([]byte, 32)
	select {
	case <-rh.Leader.Done:
		log.Printf("RandHound - done")
		bytes = <-rh.Leader.Result
	case <-time.After(time.Second * 60):
		t.Fatal("RandHound – time out")
	}
	log.Printf("RandHound - random bytes: %v\n", bytes)
}
