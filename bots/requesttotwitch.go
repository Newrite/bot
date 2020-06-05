package bots

func (bt *BotTwitch) handleRequests(cmd string) string {
	switch cmd {
	case "uptime":
		return bt.handleApiRequest("Имяюзверя", "xandr_sh", "", cmd)
	}
	return cmd
}
