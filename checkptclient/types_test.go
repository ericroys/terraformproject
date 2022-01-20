package checkptclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoginResponse(t *testing.T) {
	data := `{
		"uid" : "b78562f2-5dfe-4f5a-b9f9-9f9a04b0e1cb",
		"sid" : "xws4UEdyXg2JHwMbecEIdUJ2Wzo_fOA9AwdIDVx4uek",
		"url" : "https://r8020sms.seamlessti.net:443/web_api/v1.3",
		"session-timeout" : 600,
		"last-login-was-at" : {
		  "posix" : 1556899851157,
		  "iso-8601" : "2019-05-03T11:10-0500"
		},
		"disk-space-message" : "Partition /var/log has: 819 MB of free space and it's lower than required: 2000 MB\n",
		"api-server-version" : "1.3"
	  }`
	var i LoginResponse

	var m = []byte(data)
	err := json.Unmarshal(m, &i)
	if err != nil {
		t.Fatalf("failed to transform response message. %v", err)
	}
	fmt.Printf("%d", i.SessTimeout)

}
