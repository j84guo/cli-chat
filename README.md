## CLI Chat
A command line chat tool based on Git concepts.

### Features 
* local config
	- client checks for local config file, prompts signup or update on start
  	- username
  	- list of followed rooms (fetched on start)
* public and private chat rooms
	- view public rooms
	- checkout rooms
	- follow/mute rooms
	- enter private rooms with ssh key

### Phases

1. Local Config
Prompt user to signup if no local config found, else prompt to update or use
current config. If signup, read name from stdin and save to file.

