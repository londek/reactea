package displayplayer

import (
	"fmt"

	"github.com/londek/reactea/examples/dynamicRoutes/data"
)

type Props = int

func Render(playerId Props, width, height int) string {
	if player, ok := data.Players[playerId]; ok {
		return fmt.Sprintf("Name: %s. Year of birth: %d. Team: %s.\nPress ctrl+c to exit or U to go back!", player.Name, player.YearOfBirth, player.Team)
	} else {
		return fmt.Sprintf("Player with ID %d not found!\nPress ctrl+c to exit or U to go back!", playerId)
	}
}
