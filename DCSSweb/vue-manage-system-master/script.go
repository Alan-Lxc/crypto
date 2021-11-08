package main

import "fmt"

type activity struct {
	content   string
	timestamp string
}

func New() {

}

func main() {
	//前20位是位置信息
	str := "2021/10/15 01:58:42 Node 1 new done,ip = 127.0.0.1:10001\n2021/10/15 01:58:43 [Node 1] Now Get start Phase1\n2021/10/15 01:58:43 [Node 1] send point message to [Node 2]\n2021/10/15 01:58:43 [Node 1] send point message to [Node 3]\n2021/10/15 01:58:43 [Node 1] send point message to [Node 4]\n2021/10/15 01:58:43 [Node 1] send point message to [Node 5]\n2021/10/15 01:58:43 Phase 1 :[Node 1] receive point message from [Node 2]\n2021/10/15 01:58:43 Phase 1 :[Node 1] receive point message from [Node 3]\n2021/10/15 01:58:43 Phase 1 :[Node 1] receive point message from [Node 4]\n2021/10/15 01:58:43 Phase 1 :[Node 1] receive point message from [Node 5]\n2021/10/15 01:58:43 [Node 1] read bulletinboard in phase 1\n2021/10/15 01:58:43 [Node 1] receive zero message from [Node 2] in phase 2\n2021/10/15 01:58:43 [Node 1] send message to [Node 2] in phase 2\n2021/10/15 01:58:43 [Node 1] send message to [Node 3] in phase 2\n2021/10/15 01:58:43 [Node 1] send message to [Node 4] in phase 2\n2021/10/15 01:58:43 [Node 1] send message to [Node 5] in phase 2\n2021/10/15 01:58:43 [Node 1] receive zero message from [Node 3] in phase 2\n2021/10/15 01:58:43 [Node 1] receive zero message from [Node 4] in phase 2\n2021/10/15 01:58:43 [Node 1] receive zero message from [Node 5] in phase 2\n2021/10/15 01:58:43 1 has finish _0ShareSum\n2021/10/15 01:58:43 [Node %!d(MISSING)] start verification in phase 2\n2021/10/15 01:58:43 [Node 1] read bulletinboard in phase 2\n2021/10/15 01:58:43 [Node 1] exponentSum: [919527479152314492520017575367191592683263256082405632185091027972819030997943823490183648180157342586948451807488935955194244950067545680105211023924584, 4488708331611103817365843553572127828114720293967289052761970848310698238665156340739995637022331696132887341554588441642277309627493593404039601398517850]\n2021/10/15 01:58:44 [Node 1] receive point message from [Node 4] in phase3\n2021/10/15 01:58:44 node 1 send point message to node 2 in phase 3\n2021/10/15 01:58:44 node 1 send point message to node 3 in phase 3\n2021/10/15 01:58:44 node 1 send point message to node 4 in phase 3\n2021/10/15 01:58:44 node 1 send point message to node 5 in phase 3\n2021/10/15 01:58:44 [Node 1] receive point message from [Node 5] in phase3\n2021/10/15 01:58:44 [Node 1] receive point message from [Node 2] in phase3\n2021/10/15 01:58:44 [Node 1] receive point message from [Node 3] in phase3\n2021/10/15 01:58:44 [Node 1] has finish sharePhase3\n2021/10/15 01:58:44 [Node 1] write bulletinboard in phase 3\n"
	item := activity{
		timestamp: "2021/10/15 01:58:42",
		content:   "Node 1 new done,ip = 127.0.0.1:10001",
	}

	activities := []activity{item}
	for i := 0; i < len(str) && i < len(str); i++ {
		if str[i] == '\n' {
			j := i + 1
			if j < len(str) {
				for j = i + 1; str[j] != '\n' && j < len(str); j++ {

				}
				item = activity{
					timestamp: str[i+1 : i+20],
					content:   str[i+20 : j],
				}
				activities = append(activities, item)
			}
		}
	}
	for i := 0; i < len(activities); i++ {
		fmt.Println("	content: \"" + activities[i].content + "\",")
		fmt.Println("	timestamp: \"" + activities[i].timestamp + "\"")
		fmt.Println("},{")

	}
}
