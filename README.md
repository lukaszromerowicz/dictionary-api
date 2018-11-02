# dictionary-api

Web API providing word search from English dictionary. It finds words possible to be built with given letters.

Inspired by Biritsh game show [Countdown](https://en.wikipedia.org/wiki/Countdown_(game_show)).

## Usage

`/words?letters=cod`

Returns words possible to be built; **cod** are the letters in search.

Example response: 

```json
{
	"Count": 6,
	"Words": [
		{
			"Word": "cod",
			"Length": 3
		}, 
		{
			"Word": "do",
			"Length": 2
		}, 
		{
			"Word": "od",
			"Length": 2
		}, 
		{
			"Word": "c",
			"Length": 1
		}, 
		{
			"Word": "d",
			"Length": 1
		}, 
		{
			"Word": "o",
			"Length": 1
		}
	]
}
```
