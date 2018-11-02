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

`/words?letters=cod&size=3`

Returns a list of words that can be matched with a length of **3** or more.

Example response: 

```json
{
	"Count": 6,
	"Words": [
		{
			"Word": "cod",
			"Length": 3
		}
	]
}
```

`/words?letters=cod&limit=2`

Returns a list of words that have been limited to **2**

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
		}
	]
}
```

`/definition?word=cod`

Returns a single word that matches the search string **cod** this contains a meaning property which has an array of possible meanings

Example response: 

```json
{
	"Word": [
		{
			"Word": "cod",
			"Length": 3,
			"Meaning" : ["A husk; a pod; as, a peascod.","A small bag or pouch.","The scrotum.","A pillow or cushion.","An important edible fish (Gadus morrhua), taken in immense numbers on the northern coasts of Europe and America. It is especially abundant and large on the Grand Bank of Newfoundland. It is salted and dried in large quantities."]
		}
	]
}
```