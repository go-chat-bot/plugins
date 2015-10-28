[![Circle CI](https://circleci.com/gh/go-chat-bot/plugins.svg?style=svg)](https://circleci.com/gh/go-chat-bot/plugins)

### Active

* **gif**: Posts a random gif url from [giphy.com][giphy.com]. Try it with: **!gif cat**
* **catgif**: Posts a random cat gif url from [thecatapi.com][thecatapi.com]
* **godoc**: Searches packages in godoc.org. Try it with: **!godoc net/http**
* **puppet**: Allows you to send messages through the bot: Try it with: **!puppet say #go-bot Hello!**
* **guid**: Generates a new guid
* **crypto**: Encrypts the input data using sha1 or md5
* **encode**: Encodes a string, currently only to base64
* **decode**: Decodes a string. currently ony from base64
* **treta**: Use it to sow discord on a channel

Tip: Use `!help <command>` to obtaing more info about these commands.

### Passive (triggers)

Passive commands receive all the text sent to the bot or the channels that the bot is in and can process it and reply.

These commands differ from the active commands as they are executed for every text that the bot receives. Ex: The Chuck Norris command, replies with a Chuck Norris fact every time the words "chuck" or "norris" are mentioned on a channel.

* **url**: Detects url and gets it's title (very naive implementation, works sometimes)
* **catfacts**: Tells a random cat fact based on some cat keywords
* **jira**: Detects jira issue numbers and posts the url (necessary to configure the JIRA URL)
* **chucknorris**: Shows a random chuck norris quote every time the word "chuck" is mentioned

### Brazilian commands (pt-br)

Some commands only makes sense to brazillians:

* **megasena**: Gera um número da megasena ou mostra o último resultado
* **cotacao**: Informa a cotação atual do Dólar e Euro
* **dilma** (passivo): Diz alguma frase da Dilma quando a palavra "dilma" é citada
* **cpf**: Gera e valida CPFs
* **cnpj**: Gera e valida CNPJs

### Wish to write a new plugin?

Start with the example commands in the [example directory](https://github.com/go-chat-bot/plugins/tree/master/example).

It's dead simple, you just need to write a go function and register it on the bot.
