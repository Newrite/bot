package twitch

func (bt *BotTwitch) HandleRequests(cmd string) string {
	switch cmd {
	case "uptime":
		return bt.handleApiRequest("Имяюзверя", "xandr_sh", "", cmd)
	}
	return cmd
}
