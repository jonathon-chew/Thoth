package utils

import (
	"fmt"
	"time"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

type CommitMap map[string]int

// Step 3: Render basic ASCII heatmap
func RenderDateGraph(commits CommitMap) {
	start := time.Now().AddDate(0, 0, -365)
	for d := 0; d <= 365; d++ {
		day := start.AddDate(0, 0, d)
		if day.Day()%28 == 0 {
			fmt.Print("\t" + day.Month().String() + "\n")
		}
		key := day.Format("2006-01-02")
		count := commits[key]
		fmt.Print(heatChar(count))
	}
	fmt.Println()
}

func heatChar(count int) string {
	switch {
	case count == 0: // 0 commits in a day
		returnString, err := aphrodite.ReturnColour("Black", "§")
		if err != nil {
			return "§"
		}
		return returnString
	case count < 2: // 1-2 commits in a day
		returnString, err := aphrodite.ReturnColour("Red", "§")
		if err != nil {
			return "§"
		}
		return returnString
	case count < 5: //2-4 commits in a day
		returnString, err := aphrodite.ReturnBold("Green", "§")
		if err != nil {
			return "§"
		}
		return returnString
	case count < 10: //5-9 commits in a day
		returnString, err := aphrodite.ReturnHighIntensity("Yellow", "§")
		if err != nil {
			return "§"
		}
		return returnString
	default: // 10 or more commits in a day
		returnString, err := aphrodite.ReturnHighIntensityBackgrounds("Purple", "§")
		if err != nil {
			return "§"
		}
		return returnString
	}
}
