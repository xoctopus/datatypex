package snowflake_id_test

import (
	"net"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/gomega"

	. "github.com/sincospro/datatypes/snowflake_id"
)

func TestWorkerIDFromIP(t *testing.T) {
	wid := WorkerIDFromIP(net.ParseIP("255.255.255.255"))
	NewWithT(t).Expect(wid).To(Equal(uint32(65535)))

	wid = WorkerIDFromIP(net.ParseIP("127.0.0.1"))
	NewWithT(t).Expect(wid).To(Equal(uint32(1)))

	wid = WorkerIDFromIP(nil)
	NewWithT(t).Expect(wid).To(Equal(uint32(0)))

	t.Log(WorkerIDFromLocalIP())

	p := gomonkey.NewPatches()
	defer p.Reset()

	p = gomonkey.ApplyFuncReturn(os.Hostname, "", nil)
	p = gomonkey.ApplyFuncReturn(os.Getenv, "")

	_, err := WorkerIDFromLocalIP()
	NewWithT(t).Expect(err).NotTo(BeNil())
}

func TestRandU32N(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(RandU32N(9))
	}
}
