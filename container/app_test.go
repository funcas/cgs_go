package container

import "testing"

func TestReadConfig(t *testing.T) {

	c := readConfig("../conf/outlet.json", "outlets")
	for _, i := range c[0] {
		r := i.GetArray("acceptTransCodes")
		for _, v := range r {
			t.Log(string(v.GetStringBytes()))
		}
	}
}
