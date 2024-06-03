# whatstlunch

What's for Lunch is a handy service designed to help you search for dishes based on the ingredients you have on hand. 
It even supports ingredient recognition through AI-powered camera functionality.

My mum struggles every day to find what dish to make, so I created this service.

## Installation

```bash
git clone https://github.com/ocxide/whatstlunch.git
```

### Frontend

The frontend is crafted using Astro, SolidJS, and Tailwind. To manage dependencies, we utilize

### Backend

The server is written in `Go`. Make sure you have Go installed on your system by following the instructions. [check here](https://go.dev/doc/install).

### Build

If you have the necessary tools installed (bun and Go), you can build the entire project with a simple command:

```bash
chmod +x ./build.sh && ./build.sh
```

This command will generate a `whatstlunch` binary for the server sub-project and a public directory containing the frontend build.

### Dependencies

This service relies on `sqlite3` as its database. Ensure it is installed on your system using the following command:

```bash
sudo apt install sqlite3 # or your package manager
```

## Contributing

Feel free to open issues and PRs.
