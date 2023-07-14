<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/capcom6/service-monitor-tgbot">
    <img src="assets/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Service Monitor Telegram Bot</h3>

  <p align="center">
    Telegram bot for monitoring the availability of network services.
    <br />
    <!-- <a href="https://github.com/capcom6/service-monitor-tgbot"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/capcom6/service-monitor-tgbot">View Demo</a> -->
    ·
    <a href="https://github.com/capcom6/service-monitor-tgbot/issues">Report Bug</a>
    ·
    <a href="https://github.com/capcom6/service-monitor-tgbot/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [HTTP service monitoring example](#http-service-monitoring-example)
  - [TCP service monitoring example](#tcp-service-monitoring-example)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

<!-- ABOUT THE PROJECT -->
## About The Project

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->

Monitoring the availability of network services is an important task for any project. At the same time, it is not always necessary to deploy universal solutions like Prometheus, a fairly simple solution. It is for such cases that this bot was created.

The bot will allow you to monitor the availability of HTTP(S) and TCP services and notify Telegram about changes in their status.

The project is in the MVP stage.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Golang][Golang]][Golang-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Follow the instructions below to run a bot using Docker.

### Prerequisites

To run a bot, it is enough to have Docker or another environment for running containers.

### Installation

1. Register a new bot and get a token for it: https://core.telegram.org/bots/features#creating-a-new-bot
2. Create a [channel](https://telegram.org/tour/channels) or [group](https://telegram.org/tour/groups) to which notifications will be sent
3. Add the bot to the channel/group as an administrator with the ability to send messages
4. Copy the configuration file [config.example.yml](configs/config.example.yml) to your working directory as `config.yml`
5. Make changes to the configuration file:
    - specify the bot token;
    - specify the channel/group ID, you can find out the ID, for example, by following the link like [https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]](https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]) after adding the bot to a channel/group and finding the value of `my_chat_member.chat.id`;
    - list the services to monitor.
6. Run the docker container: `docker run -d -v "$(pwd)/config.yml:/app/config.yml:ro" --name tgbot capcom6/service-monitor-tgbot:latest`

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

### HTTP service monitoring example

```yaml
services:
  - name: Google
    initialDelaySeconds: 5
    periodSeconds: 10
    timeoutSeconds: 1
    successThreshold: 1
    failureThreshold: 3
    httpGet:
      scheme: https
      host: google.com
      path: /
      port: 443
      httpHeaders:
        - name: X-Header
          value: value
```

![HTTP Alert][http-alert]

### TCP service monitoring example

```yaml
services:
  - name: MySQL
    initialDelaySeconds: 5
    periodSeconds: 10
    timeoutSeconds: 1
    successThreshold: 1
    failureThreshold: 3
    tcpSocket:
      host: localhost
      port: 3306
```

![TCP Alert][tcp-alert]

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Add Changelog
- [x] Add the ability to change the text of messages
- [ ] Send notifications to multiple channels/groups
- [ ] Display event time in notifications
- [ ] Online/offline time count
- [ ] Active bot mode
     - [ ] Request current state of services
     - [ ] SLA report
     - [ ] The event log
- [ ] Separation of bot and monitoring service
- [ ] Dynamic list of services
- [ ] Service discovery

See the [open issues](https://github.com/capcom6/service-monitor-tgbot/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the Apache-2.0 license. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Project Link: [https://github.com/capcom6/service-monitor-tgbot](https://github.com/capcom6/service-monitor-tgbot)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
<!-- ## Acknowledgments

Use this space to list resources you find helpful and would like to give credit to. I've included a few of my favorites to kick things off!

* [Choose an Open Source License](https://choosealicense.com)
* [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
* [Malven's Flexbox Cheatsheet](https://flexbox.malven.co/)
* [Malven's Grid Cheatsheet](https://grid.malven.co/)
* [Img Shields](https://shields.io)
* [GitHub Pages](https://pages.github.com)
* [Font Awesome](https://fontawesome.com)
* [React Icons](https://react-icons.github.io/react-icons/search)

<p align="right">(<a href="#readme-top">back to top</a>)</p> -->



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[contributors-url]: https://github.com/capcom6/service-monitor-tgbot/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[forks-url]: https://github.com/capcom6/service-monitor-tgbot/network/members
[stars-shield]: https://img.shields.io/github/stars/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[stars-url]: https://github.com/capcom6/service-monitor-tgbot/stargazers
[issues-shield]: https://img.shields.io/github/issues/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[issues-url]: https://github.com/capcom6/service-monitor-tgbot/issues
[license-shield]: https://img.shields.io/github/license/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[license-url]: https://github.com/capcom6/service-monitor-tgbot/blob/master/LICENSE.txt
[product-screenshot]: assets/screenshot.png
[http-alert]: assets/http-alert.png
[tcp-alert]: assets/tcp-alert.png
[Golang]: https://img.shields.io/badge/Golang-000000?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
