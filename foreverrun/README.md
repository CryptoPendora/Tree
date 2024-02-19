# Graceful Termination Example

This project demonstrates a simple Go program that prints a continuously incrementing number prefixed with "test" and a message instructing the user to press ctrl+c to exit. The program uses graceful termination to handle SIGINT and SIGTERM signals, ensuring that the program exits cleanly.

## How It Works

- The program starts a goroutine that runs `printTestContinuously`, which prints "test" followed by a number. The number increments every second, and the line is overwritten to give the illusion of an updating message.
- The main function calls `gracefulTermination` to handle graceful termination. The program waits for a SIGINT or SIGTERM signal, upon which it prints a message indicating that the program is exiting gracefully.

## Features

- **Continuous Output**: The program prints "test x" where x is a continuously incrementing number. If the number exceeds 9999, it resets to 0 to avoid overflow.
- **Graceful Termination**: The program handles SIGINT and SIGTERM signals. When these signals are received, the program exits gracefully, ensuring that all resources are properly released.

## Usage
![foreverrun](https://github.com/godevgod/foreverrun/blob/main/foreverrun.gif?raw=true)

1. Clone the repository to your local machine.
2. Navigate to the directory containing the program.
3. Run the program using the command:

   ```bash
   go run main.go

### Support
If you find this project helpful and would like to support my work, feel free to buy me a drink:

# Your support is greatly appreciated!

![Bitcoin QR Code](https://github.com/godevgod/foreverrun/blob/main/1CbE3SsUcvJWZ2YNaDwUj9AQtT8k4AGmLe.png?raw=true)
### Bitcoin:1CbE3SsUcvJWZ2YNaDwUj9AQtT8k4AGmLe
```bash
   1CbE3SsUcvJWZ2YNaDwUj9AQtT8k4AGmLe


