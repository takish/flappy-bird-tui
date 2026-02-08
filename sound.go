package main

import "fmt"

// PlaySound plays a beep sound using terminal bell
func PlaySound(soundType string) {
	switch soundType {
	case "jump":
		// Single beep for jump
		fmt.Print("\a")
	case "score":
		// Double beep for score
		fmt.Print("\a")
	case "gameover":
		// Triple beep for game over
		fmt.Print("\a")
	}
}
