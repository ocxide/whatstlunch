# whatstlunch

whatstlunch is a handy service designed to help you search for dishes based on the ingredients you have on hand. 
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

## Usage

After building the project, you will need to populate the database with some data. For example, run:

```bash
whatstlunch load https://www.recetasgratis.net/busqueda/q/plato%20principal/meal_id/3
```

This will load into database the data from recetasgratis.net, removing the previous data.
You can customize the filters in the webside in order to load specific data.

Currently, the only informantion source supported is recetasgratis.net, more comming soon.

Then, launch the service to http with:

```bash
whatstlunch launch
```

By default this will list on host 127.0.0.1 and port 3456. Visit http://127.0.0.1:3456/ to see the client.

### AI ingredients Recognition

This service can recognize ingredients from images taken with your phone or a selected image. For this to work, you will need to specify an
"chatgpt-like" API to interanct.

The recommended AI service is  `llava:7b` model of [ollama](https://www.ollama.com/).
Specifing an API_KEY is not currently supported yet, so, this will probably not work with chatgpt.

Create a config file in the executable directory:

```toml
// config.toml

[ai]
model="llava:7b"
api_url="http://127.0.0.1:11434/api"
```

See the default config in the `whatstlunch-server/config.toml` file.

## Contributing

Feel free to open issues and PRs.
