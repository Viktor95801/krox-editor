import std.conv;
import std.file;
import std.json;
import std.random;
import std.stdio;

import bindbc.sdl;
import bindbc.loader;

enum LogLevel
{
	DEBUG,
	INFO,
	WARNING,
	ERROR,
	SUCCESS,
}

void log(int level, string arg)
{
	import std.stdio;

	switch (level)
	{
		debug
		{
	case LogLevel.DEBUG:
			writeln("\033[1;34m[DEBUG]:\033[0m", arg);
			break;
		}
	case LogLevel.INFO:
		writeln("\033[1;34m[INFO]: \033[0m", arg);
		break;
	case LogLevel.WARNING:
		writeln("\033[33m[WARNING]: \033[0m", arg);
		break;
	case LogLevel.ERROR:
		writeln("\033[31m[ERROR]: \033[0m", arg);
		break;
	case LogLevel.SUCCESS:
		writeln("\033[1;32m[SUCCESS]: \033[0m", arg);
		break;
	default:
		if (level == LogLevel.DEBUG)
		{
			break;
		}
		assert(0);
	}
}

// structs
struct TextArea
{
	// SDL
	TTF_Font* font;

	// Krox
	string[] lines;
	int cursor_pos_max = 0;
	int cursor_pos = 0;
	int cursor_line = 0;
}

// Krox
int main(string[] args)
{
	//version (Windows)
	//	setCustomLoaderSearchPath(".\\libs\\SDL3.dll");
	LoadMsg ret = loadSDL();
	if (ret != LoadMsg.success)
	{
		foreach (info; errors)
		{
			log(LogLevel.ERROR, to!string(info.message));
		}
		switch (ret)
		{
		case LoadMsg.noLibrary:
			log(LogLevel.ERROR, "Couldn't find/load sdl3.dll");
			break;
		case LoadMsg.badLibrary:
			string msg = "Your SDL version is not supported, please upgrade to 3.2.0+.";
			log(LogLevel.ERROR, msg);
			break;
		default:
			assert(0);
		}
	}
	else
	{
		log(LogLevel.SUCCESS, "Loaded successfully sdl3.dll");
	}

	SDL_Init(SDL_INIT_VIDEO);
	TTF_Init();

	string continuation;
	{
		string jsonContents = readText("dist/bin/krox/continuations.json");

		JSONValue parsed = parseJSON(jsonContents);

		auto data = parsed["continuations"].array;

		continuation = data[uniform(0, data.length)].str;
	}
	SDL_Window* window = SDL_CreateWindow(("Krox: " ~ continuation).ptr,
		800, 600, 0);

	if (window is null)
	{
		log(LogLevel.ERROR, "Couldn't create a SDL window, error message: ");
		log(LogLevel.ERROR, to!string(SDL_GetError()));
		SDL_Quit();
	}

	// main event loop
	SDL_Event event;
	auto running = true;
	while (running)
	{
		while (SDL_PollEvent(&event) != 0)
		{
			if (event.type == SDL_EVENT_QUIT)
			{
				running = false;
			}
		}
		SDL_Delay(16); // Sleep to reduce CPU usage
	}

	SDL_DestroyWindow(window);
	SDL_Quit();
	log(LogLevel.INFO, "Exit with success");

	return 0;
}
