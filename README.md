# bedrock-server-downloads

An automated system that tracks and stores data about the Minecraft Bedrock Dedicated Server. The system searches for and records the latest download links for both release and preview versions.

## Structure

* `bedrock-server-downloads.json`: JSON file storing download links for Minecraft Bedrock servers for both Windows and Linux
* `main.go`: Main code used to fetch data and update the JSON file

## Operation

This system works automatically through GitHub Actions, which runs a script to check for new versions of Minecraft Bedrock Server from the official Minecraft website daily. If a new version is found, the system updates the JSON file and commits the changes to this repository.

Features include:
- Tracking both release and preview versions
- Supporting both Windows and Linux platforms
- Maintaining a history of download links for all versions

## Usage

You can use our JSON file API to retrieve the latest download links or specific versions of Minecraft Bedrock Server.

```
https://raw.githubusercontent.com/kittizz/bedrock-server-downloads/main/bedrock-server-downloads.json
```

## Contributing

If you would like to contribute or find any errors, please fork this repository and create a pull request or open an issue for discussion.

## Acknowledgements

Special thanks to [EndstoneMC/bedrock-server-data](https://github.com/EndstoneMC/bedrock-server-data) for the concept and examples that helped in developing this system.

## License

This project is under the MIT License. Please see the LICENSE file for more details.