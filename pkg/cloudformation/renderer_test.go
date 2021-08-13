package cloudformation

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	ren := &Renderer{
		Name: "TestTask",
		Params: map[string]interface{}{
			"EnvType": "production",
		},
	}
	cfm := "../../test/task.yaml"
	indent, min, err := ren.Render(cfm)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%s\n", string(indent))
		m := make(map[string]interface{})
		if err = json.Unmarshal(min, &m); err != nil {
			t.Error(err)
		} else {
			defs := m["containerDefinitions"]
			for _, v := range defs.([]interface{}) {
				container := v.(map[string]interface{})
				if container["name"] != "my-app" {
					continue
				}
				if container["image"] != "xxxxxxxxxxxx.dkr.ecr.xxxxxx-1.amazonaws.com/test-image:latest" {
					t.Errorf("image=%v", container["image"])
					return
				}
				if container["cpu"] != float64(256) {
					t.Errorf("cpu=%v", container["cpu"])
					return
				}
				if container["memory"] != float64(512) {
					t.Errorf("memory=%v", container["memory"])
					return
				}
				for _, val := range container["environment"].([]interface{}) {
					env := val.(map[string]interface{})
					switch env["name"] {
					case "LOG_LEVEL":
						if env["value"] != "info" {
							t.Errorf("LOG_LEVEL=%v", env["value"])
						}
					case "GIN_MODE":
						if env["value"] != "release" {
							t.Errorf("GIN_MODE=%v", env["value"])
						}
					}
				}
			}
		}
	}
}
