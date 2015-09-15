package fixtures
import (
	"io/ioutil"
	"bufio"
	"os"
)

// Make temporary BLOB
func MakeBLOB(size int64) (name string, err error)  {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer f.Close()
	name = f.Name()

	var i int64
	var j uint8
	buf := bufio.NewWriter(f)

	for i=0; i < size; i++ {
		if _, err = buf.Write([]byte{j}); err != nil {
			return
		}
		j++
		if j > 126 {
			j = 0
		}
	}
	err = buf.Flush()
	return
}

func KillBLOB(name string) (err error) {
	return os.Remove(name)
}
