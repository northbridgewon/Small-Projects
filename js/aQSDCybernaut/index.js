/**
 * @fileOverview A Discord Bot written in JavaScript!
 * @author Ak Yair Lin Cortek (northbrigdewon.dev) + Gemini (gemini.google.com)
 * @version 1.0.0
 * @license MIT
 *
 * This script provides utility function.
 *
 * For the full license text, please see the LICENSE.md file
 * or visit https://opensource.org/licenses/MIT
 */

(function() {
    'use strict';

// This file is now the conductor of our bot orchestra!

const fs = require('node:fs'); // Node.js file system module, for reading command files
const path = require('node:path'); // Node.js path module, for constructing file paths
const { Client, GatewayIntentBits, Events, Collection } = require('discord.js');

// Load configuration
// We'll try to load config.json. If it's not there, we'll guide the user.
let config;
try {
    const configPath = path.join(__dirname, 'config.json');
    config = require(configPath); // require() can directly read JSON files
} catch (error) {
    console.error("💀 Whoopsie! config.json is missing or malformed.");
    console.error("Please create a config.json file in the root directory with your BOT_TOKEN.");
    console.error("You can copy config.example.json to config.json and fill in your details.");
    process.exit(1); // Exit if config is not found
}

const token = config.BOT_TOKEN;

if (!token) {
    console.error("💔 Major heartbreak! The BOT_TOKEN is missing from your config.json. Can't start without it!");
    process.exit(1);
}

// Create a new Client instance
const client = new Client({
    intents: [
        GatewayIntentBits.Guilds,
        GatewayIntentBits.GuildMessages,
        GatewayIntentBits.MessageContent // Remember to enable this in the Developer Portal!
    ]
});

// Storing our commands
// We use a discord.js Collection to store commands. It's like an enhanced Map.
client.commands = new Collection();

// Dynamically load command files
// This is where the modular magic happens! ✨
const commandsPath = path.join(__dirname, 'commands'); // Path to the 'commands' directory
try {
    const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js')); // Get all .js files

    for (const file of commandFiles) {
        const filePath = path.join(commandsPath, file);
        const command = require(filePath); // Load the command module

        // Set a new item in the Collection with the key as the command name and the value as the exported module
        if ('name' in command && 'execute' in command) {
            client.commands.set(command.name, command);
            console.log(`✅ Command loaded: ${command.name} from ${file}`);
        } else {
            console.warn(`[WARNING] The command at ${filePath} is missing a required "name" or "execute" property. Skippin' it!`);
        }
    }
} catch (error) {
    console.error(`🚫 Oh noes! Could not read the commands directory at ${commandsPath}:`, error);
    console.error("Make sure you have a 'commands' folder with your command files in it.");
    // Depending on how critical commands are, you might want to process.exit(1) here.
}


// ClientReady event - runs once when the bot is ready
client.once(Events.ClientReady, readyClient => {
    console.log(`🚀 Logged in as ${readyClient.user.tag}! The modular bot is ready to party!`);
    readyClient.user.setActivity("managing modular commands");
});

// MessageCreate event - runs every time a message is created
client.on(Events.MessageCreate, async message => {
    if (message.author.bot) return; // Ignore messages from bots (including ourself!)

    const prefix = config.PREFIX || "!"; // Use prefix from config or default to "!"

    if (!message.content.startsWith(prefix)) return; // Only process messages with our prefix

    // Parse the command and arguments
    const args = message.content.slice(prefix.length).trim().split(/ +/);
    const commandName = args.shift().toLowerCase();

    const command = client.commands.get(commandName) || client.commands.find(cmd => cmd.aliases && cmd.aliases.includes(commandName));

    if (!command) {
        // Optional: reply if the command doesn't exist
        await message.reply("Hmm, I don't know that command! 🤔");
        console.log(`Command not found: ${commandName}`);
        return;
    }

    try {
        // Execute the command!
        await command.execute(message, args, client, config); // Pass client and config if commands need them
        console.log(`Executed command '${command.name}' for ${message.author.tag}`);
    } catch (error) {
        console.error(`💥 Error executing command '${command.name}':`, error);
        await message.reply('Yikes! There was an error trying to execute that command! 😵‍💫').catch(console.error);
    }
});

// Log in to Discord with your client's token
client.login(token)
    .catch(error => {
        console.error("💀 Catastrophic failure to login! Double-check that BOT_TOKEN in config.json.", error);
    });

// Enhanced error handling
process.on('unhandledRejection', error => {
    console.error('Unhandled promise rejection:', error);
    // Consider more sophisticated logging or alerting here for a production bot
});

process.on('uncaughtException', error => {
    console.error('Uncaught exception:', error);
    // It's often recommended to restart the process on an uncaught exception,
    // as the application state might be corrupted.
    // process.exit(1);
});

console.log("尝试启动机器人... (Attempting to start the bot...)");

})(); 

/**
 * @license
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */