# Weather Fetcher CLI

A simple yet powerful command-line interface (CLI) tool built with Go to fetch and display current weather information and 5-day forecasts from anywhere in the world.

---

## Features

-   **Current Weather:** Get the real-time temperature, condition, and "feels like" temperature for any city.
-   **5-Day Forecast:** View a daily weather summary for the next five days.
-   **Configurable Units:** Display temperatures in Celsius (`metric`) or Fahrenheit (`imperial`).
-   **Smart Caching:** Avoids excessive API calls by caching recent results for 10 minutes, providing faster responses.
-   **Secure API Key Management:** Reads your API key from an environment variable to keep it out of your source code.

## Installation & Setup

To get started with the Weather Fetcher CLI, you'll need to have Go installed on your system.

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/your-username/weather-cli.git
    cd weather-cli
    ```

2.  **Get an API Key:**
    This tool uses the [OpenWeatherMap API](https://openweathermap.org/api). You will need to sign up for a free account to get your personal API key.

3.  **Set the Environment Variable:**
    Make your API key available to the application by setting it as an environment variable.

    *   On macOS and Linux:
        ```bash
        export WEATHER_API_KEY='your_api_key_here'
        ```
    *   On Windows (PowerShell):
        ```powershell
        $env:WEATHER_API_KEY='your_api_key_here'
        ```
    > **Note:** To make this variable permanent, you should add this line to your shell's startup file (e.g., `.bashrc`, `.zshrc`, or your PowerShell Profile).

4.  **Build the Executable:**
    From the project's root directory, run the `go build` command. This will compile the source code into a single executable file named `weather-cli`.
    ```bash
    go build
    ```

## Usage

The executable (`./weather-cli` on Linux/macOS or `weather-cli.exe` on Windows) accepts several flags to customize your request.

**The `-city` flag is always required.**

---

### Get Current Weather

To get the current weather, simply provide a city name. The default unit is Celsius (`metric`).

```bash
./weather-cli -city="London"