<!-- Improved compatibility of back to top link: See: https://github.com/capcom6/tgbot-service-monitor/Best-README-Template/pull/73 -->
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
  <a href="https://github.com/capcom6/tgbot-service-monitor">
    <img src="assets/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Service Monitor Telegram Bot</h3>

  <p align="center">
    Телеграм-бот для мониторинга доступности сетевых сервисов.
    <!-- <br />
    <a href="https://github.com/capcom6/tgbot-service-monitor/Best-README-Template"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/capcom6/tgbot-service-monitor/Best-README-Template">View Demo</a>
    ·
    <a href="https://github.com/capcom6/tgbot-service-monitor/Best-README-Template/issues">Report Bug</a>
    ·
    <a href="https://github.com/capcom6/tgbot-service-monitor/Best-README-Template/issues">Request Feature</a> -->
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

[![Product Name Screen Shot][product-screenshot]](https://example.com)

Мониторинг доступности сетевых сервисов - важная задача для любого проекта. При этом не всегда есть необходимость разворачивать универсальные решения типа Prometheus, а достаточно простого решения. Именно для таких случаев и был создан данный бот.

Бот позволит мониторить доступность HTTP(S) и TCP сервисов и уведомлять об изменении их состояния в Telegram.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Golang][Golang]][Golang-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Для запуска бота с использованием Docker следуйте инструкциям ниже.

### Prerequisites

Для запуска бота достаточно наличия Docker или иного окружения для запуска контейнеров.

### Installation

1. Зарегистрируйте нового бота и получите токен для него: https://core.telegram.org/bots/features#creating-a-new-bot
2. Создайте [канал](https://telegram.org/tour/channels) или [группу](https://telegram.org/tour/groups) в которую будут отправляться уведомления
3. Добавьте бота в канал/группу в качестве администратора с возможностью отправки сообщений
3. Скопируйте конфигурационный файл [config.example.yml](configs/config.example.yml) в рабочую директорию.
4. Внесите изменения в конфигурационный файл:
    - укажите токен бота;
    - укажите ИД канала/группы, узнать ИД можно, например, перейдя по ссылке вида [](https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]) после добавления бота в канал/группу и найдя значение `my_chat_member.chat.id`;
    - перечислите сервисы для мониторинга.
5. Запустите docker-контейнер: `docker run -d -v "$(pwd)/config.yml:/app/config.yml:ro" --name tgbot capcom6/service-monitor-tgbot:latest`

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Add Changelog
- [x] Add back to top links
- [ ] Add Additional Templates w/ Examples
- [ ] Add "components" document to easily copy & paste sections of the readme
- [ ] Multi-language Support
    - [ ] Chinese
    - [ ] Spanish

See the [open issues](https://github.com/capcom6/tgbot-service-monitor/Best-README-Template/issues) for a full list of proposed features (and known issues).

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

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Your Name - [@your_twitter](https://twitter.com/your_username) - email@example.com

Project Link: [https://github.com/your_username/repo_name](https://github.com/your_username/repo_name)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

Use this space to list resources you find helpful and would like to give credit to. I've included a few of my favorites to kick things off!

* [Choose an Open Source License](https://choosealicense.com)
* [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
* [Malven's Flexbox Cheatsheet](https://flexbox.malven.co/)
* [Malven's Grid Cheatsheet](https://grid.malven.co/)
* [Img Shields](https://shields.io)
* [GitHub Pages](https://pages.github.com)
* [Font Awesome](https://fontawesome.com)
* [React Icons](https://react-icons.github.io/react-icons/search)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/capcom6/tgbot-service-monitor.svg?style=for-the-badge
[contributors-url]: https://github.com/capcom6/tgbot-service-monitor/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/capcom6/tgbot-service-monitor.svg?style=for-the-badge
[forks-url]: https://github.com/capcom6/tgbot-service-monitor/network/members
[stars-shield]: https://img.shields.io/github/stars/capcom6/tgbot-service-monitor.svg?style=for-the-badge
[stars-url]: https://github.com/capcom6/tgbot-service-monitor/stargazers
[issues-shield]: https://img.shields.io/github/issues/capcom6/tgbot-service-monitor.svg?style=for-the-badge
[issues-url]: https://github.com/capcom6/tgbot-service-monitor/issues
[license-shield]: https://img.shields.io/github/license/capcom6/tgbot-service-monitor.svg?style=for-the-badge
[license-url]: https://github.com/capcom6/tgbot-service-monitor/blob/master/LICENSE.txt
[product-screenshot]: images/screenshot.png
[Golang]: https://img.shields.io/badge/Golang-000000?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
[Golang]: https://img.shields.io/badge/Golang-000000?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
