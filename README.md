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
    Телеграм-бот для мониторинга доступности сетевых сервисов.
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

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->

Мониторинг доступности сетевых сервисов - важная задача для любого проекта. При этом не всегда есть необходимость разворачивать универсальные решения типа Prometheus, а достаточно простого решения. Именно для таких случаев и был создан данный бот.

Бот позволит мониторить доступность HTTP(S) и TCP сервисов и уведомлять об изменении их состояния в Telegram.

Проект находится в стадии MVP.

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
3. Скопируйте конфигурационный файл [config.example.yml](configs/config.example.yml) в рабочую директорию под именем `config.yml`.
4. Внесите изменения в конфигурационный файл:
    - укажите токен бота;
    - укажите ИД канала/группы, узнать ИД можно, например, перейдя по ссылке вида [https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]](https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]) после добавления бота в канал/группу и найдя значение `my_chat_member.chat.id`;
    - перечислите сервисы для мониторинга.
5. Запустите docker-контейнер: `docker run -d -v "$(pwd)/config.yml:/app/config.yml:ro" --name tgbot capcom6/service-monitor-tgbot:latest`

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

### Пример мониторинга HTTP-сервиса

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

### Пример мониторинга TCP-сервиса

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

- [ ] Добавить Changelog
- [ ] Добавить возможность изменять текст сообщений
- [ ] Отправка уведомлений в несколько каналов/групп
- [ ] Отражение времени события в уведомлениях
- [ ] Подсчет времени онлайн/оффлайн
- [ ] Активный режим бота
    - [ ] Запрос текущего состояния сервисов
    - [ ] Отчет по SLA
    - [ ] Журнал событий
- [ ] Разделение бота и сервиса мониторинга
- [ ] Динамический список сервисов
- [ ] Auto-discovery (в первую очередь через Docker)

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
