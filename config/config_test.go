package config

import "testing"

func TestDrop(t *testing.T) {
	r := Route{
		Upstream: []Upstream{
			{Uri: "hello",
				Weight: 0.5,
			},
		},
	}
	result := r.Drop()
	if result.Uri != "hello" {
		t.Errorf("Single route wasn't selectedd %v", result)
	}
}
func TestDropWithTwo(t *testing.T) {
	r := Route{
		Upstream: []Upstream{
			{Uri: "path_1",
				Weight: 0.5,
			},
			{Uri: "path_2",
				Weight: 0.5,
			},
		},
	}
	resultCount1 := 0
	resultCount2 := 0
	for  i := 1; i <= 1000; i++  {
		result := r.Drop()
		if result.Uri == "path_1" {
			resultCount1++
		}
		if result.Uri == "path_2" {
			resultCount2++
		}
	}
	if (resultCount1 >= 510 && resultCount1 <=490) || (resultCount2 >= 510 && resultCount2 <=490){
		t.Errorf("Bad distribution %d to %d", resultCount1, resultCount2)
	}
}
