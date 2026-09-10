<p align="center">
  <img src="frontend/public/brand/tokenhub-logo.png" alt="TokenHub" width="96" />
</p>

<h1 align="center">TokenHub</h1>

<p align="center">
  TokenHub — это корпоративная инфраструктура управления токенами (Token Governance) для AI: маршрутизация моделей, контроль доступа, оптимизация затрат, сверка счетов с провайдерами и управляемый доступ ко всем вышестоящим поставщикам AI-моделей.
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-blue.svg" alt="License" /></a>
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" alt="Docker Compose" />
  <img src="https://img.shields.io/badge/OpenAI-Compatible-10A37F" alt="OpenAI Compatible" />
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README.zh-CN.md">简体中文</a> | <a href="README.ja.md">日本語</a> | Русский
</p>

## Корпоративное управление токенами (Enterprise Token Governance)

TokenHub предоставляет предприятиям уровень управления жизненным циклом AI-моделей: от подключения провайдеров и проектных ключей до политик маршрутизации, атрибуции использования, контроля бюджетов и сверки счетов.

Ключевая проблема возникает, когда каждая команда, приложение и рабочий процесс начинают потреблять всё больше моделей и токенов. TokenHub помещает механизмы управления перед каждым вызовом модели:

- **Маршрутизация моделей:** выбор подходящей модели в зависимости от сценария, стоимости, производительности, состояния работоспособности и политик резервного переключения (fallback).
- **Управление правами доступа:** распределение и контроль доступа к токенам между сотрудниками, командами, проектами и приложениями.
- **Экономия токенов:** снижение затрат на AI с помощью кэширования, подбора моделей, квот и оптимизированных стратегий вызова.
- **Сверка счетов с провайдерами:** сопоставление внутреннего использования со счетами поставщиков, позволяющее финансовым, платформенным и бизнес-командам обосновать реальные расходы на AI.

## Почему именно TokenHub

Многие open-source AI-шлюзы сосредоточены на веерной рассылке запросов к провайдерам (один эндпоинт для вызова множества вышестоящих сервисов). Это помогает разработчикам подключать модели, но само по себе не решает задач корпоративной эксплуатации. TokenHub построен вокруг недостающего уровня управления:

- Распределение токенов управляется на уровне проектов и команд вместо копирования исходных API-ключей провайдеров в каждое приложение.
- Доступ к моделям, маршрутизация и резервное переключение являются политиками, которые администраторы могут изменять без модификации клиентского кода.
- Счета и история запросов сопоставляются с внутренними владельцами, позволяя финансовым, платформенным и бизнес-командам прозрачно анализировать расходы на AI.
- Рабочие пространства пользователя, тимлида и администратора разделяют повседневное использование, согласование, атрибуцию затрат и эксплуатацию платформы в соответствии с зонами ответственности.

## Скриншоты

<p align="center">
  <img src="docs/assets/screenshots/tokenhub-tour.webp" alt="Обзор возможностей TokenHub: вход, сводка, документация API, каналы провайдеров, каталог моделей, политики маршрутизации, аналитика использования и системные настройки" width="100%">
</p>

## Создано вокруг трёх ролей

TokenHub разделяет повседневное использование моделей, управление командами и администрирование платформы, чтобы корпоративные пользователи видели только те рабочие процессы, которые соответствуют их обязанностям.

| Роль | Фокус рабочего пространства | Руководство |
| --- | --- | --- |
| Пользователь | Поиск доступных моделей, создание API-ключей с областью видимости проекта, вызов API моделей и просмотр личной статистики использования | [Руководство пользователя](docs/ru/user-guide.md) |
| Тимлид | Управление пространствами проектов, участниками, ключами проектов, командными отчётами и атрибуцией затрат проекта | [Руководство тимлида](docs/ru/team-leader-guide.md) |
| Администратор | Настройка провайдеров, каталога моделей, политик маршрутизации, источников идентификации, RBAC, аудита и контроля затрат | [Руководство администратора](docs/ru/administrator-guide.md) |

## Возможности платформы

- **Управление ключами на уровне проектов** с привязкой к командам, правами участников, квотами и контролем параллелизма.
- **Каталог моделей и политики маршрутизации** с приоритетами, весами, порядком переключения при сбоях, выбором с учётом сценария и диагностикой работоспособности маршрутов.
- **Аналитика использования и журналы запросов**, атрибутированные к пользователю, проекту, команде, модели и центру затрат.
- **Контроль затрат** для бюджетов токенов, сравнения расходов по провайдерам, выбора моделей и будущей экономии за счёт кэширования.
- **Настройка источников идентификации** для корпоративного входа через OAuth/OIDC, а также ролевая модель (RBAC) и журналы аудита.
- **OpenAI-совместимые API моделей:** `/v1/chat/completions`, `/v1/responses`, `/v1/embeddings`; API Anthropic Messages: `/v1/messages`, `/v1/messages/count_tokens`.
- **OpenAI-совместимая генерация изображений и редактирование по образцу** через `/v1/images/generations` и `/v1/images/edits` с асинхронными задачами и хранением изображений на стороне сервера.
- **Лаконичная консоль** с компактной навигацией с учётом ролей, глобальным поиском, светлой/тёмной темой и документацией API с разделением экрана на навигацию и детали.
- **Приватное развёртывание с приоритетом SQLite** с нативными вариантами запуска через systemd и Docker Compose.
- **Поддержка многоинстансного развёртывания на PostgreSQL:** общее состояние через удалённую базу данных PostgreSQL, горизонтальное масштабирование реплик фронтенда и бэкенда, а также настройка пулов соединений. См. [руководство по развёртыванию](docs/deployment.md) и [руководство по настройке PostgreSQL](docs/postgresql-setup.md).
- **Переключение языка консоли** между английским, китайским, японским и русским.

## Экосистема провайдеров (Provider Ecosystem)

Поддержка провайдеров — это граница интеграции TokenHub. Хостинговые API, каналы подписок, локальные модели и кастомные апстримы подключаются через абстракцию Provider, поэтому единые корпоративные политики могут управлять любым маршрутом вызова моделей.

После настройки маршрутизации, разрешений, экономии токенов, атрибуции, аудита и сверки счетов, TokenHub подключает эти контролируемые рабочие процессы к OpenAI, Azure OpenAI, Anthropic, Gemini, DeepSeek, Qwen, подпискам Codex, локальным моделям и кастомным OpenAI-совместимым апстримам.

TokenHub включает нативные адаптеры Provider для OpenAI, Azure OpenAI, Anthropic, Gemini, DeepSeek, Qwen, подписок Codex и локальных моделей, а также каталог из более чем 150 шаблонов провайдеров. Популярные интеграции включают:

<p align="center">
  <img src="docs/assets/provider-showcase.svg" alt="Популярные интеграции провайдеров TokenHub среди коммерческих, подписных, локальных и кастомных апстримов." width="100%">
</p>

Шаблоны провайдеров используют соответствующий нативный адаптер, если он доступен; в противном случае они подключаются через OpenAI-совместимый эндпоинт. Набор моделей и возможностей зависит от вышестоящего сервиса и учётной записи, в то время как корпоративная политика остаётся централизованной в TokenHub.

## Быстрый старт

Нативный релиз на хосте Linux под управлением systemd:

```bash
curl -fsSL https://raw.githubusercontent.com/astaxie/TokenHub/main/deploy/native/install.sh \
  -o /tmp/tokenhub-install.sh
sudo bash /tmp/tokenhub-install.sh install
```

Docker Compose из склонированного репозитория:

```bash
cp deploy/.env.example deploy/.env
# Замените все значения change-me в deploy/.env на надёжные секреты.
./deploy/install.sh
```

Откройте в браузере:

- Консоль администратора: `http://localhost:3000`
- API бэкенда: `http://localhost:8080`
- Проверка работоспособности (Health check): `http://localhost:8080/healthz`

Данные для первоначального входа администратора:

- Имя пользователя: `admin`
- Пароль при нативной установке: выводится установщиком один раз при инсталляции
- Пароль для Docker: значение переменной `TOKENHUB_BOOTSTRAP_ADMIN_PASSWORD`

Нативный установщик проверяет контрольные суммы релизов, устанавливает службу systemd и предоставляет прямое управление обновлениями, откатом версий и перезапуском прямо из панели версий. Стандартное развёртывание в Docker запускает бэкенд и консоль в одном управляемом контейнере и обеспечивает те же элементы прямого управления без проброса Docker-сокета. Пакеты релизов сохраняются в томе `tokenhub-releases`, поэтому стандартные перезапуски и пересоздания контейнеров сохраняют применённые через панель обновления. Многоинстансные конфигурации Docker используют управляемые оператором обновления через Compose, обеспечивая одновременное обновление всех реплик. Подробнее о каждом режиме см. в [руководстве по развёртыванию](docs/ru/deployment.md).

## Документация

- [Главная страница документации](docs/ru/README.md)
- [Архитектура](docs/ru/architecture.md)
- [Руководство пользователя](docs/ru/user-guide.md)
- [Руководство тимлида](docs/ru/team-leader-guide.md)
- [Руководство администратора](docs/ru/administrator-guide.md)
- [Руководство по участию в разработке](CONTRIBUTING.ru.md)

## Участники проекта (Contributors)

TokenHub развивается благодаря отзывам о продукте, интеграциям шлюзов, документации, тестам и постоянной заботе специалистов, которые эксплуатируют его в реальных корпоративных средах.

<!-- readme: contributors -start -->

<table>
  <tr>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/astaxie">
        <img src="https://avatars.githubusercontent.com/u/233907?v=4" width="80px" alt="astaxie" />
        <br /><sub><b>astaxie</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/deepjerry-ai">
        <img src="https://avatars.githubusercontent.com/u/262369278?v=4" width="80px" alt="deepjerry-ai" />
        <br /><sub><b>deepjerry-ai</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/legendtkl">
        <img src="https://avatars.githubusercontent.com/u/2370761?v=4" width="80px" alt="legendtkl" />
        <br /><sub><b>legendtkl</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/Mr0bean">
        <img src="https://avatars.githubusercontent.com/u/19573968?v=4" width="80px" alt="Mr0bean" />
        <br /><sub><b>Mr0bean</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/cngump">
        <img src="https://avatars.githubusercontent.com/u/108251?v=4" width="80px" alt="cngump" />
        <br /><sub><b>cngump</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/bailu-ZZ">
        <img src="https://avatars.githubusercontent.com/u/311096537?v=4" width="80px" alt="bailu-ZZ" />
        <br /><sub><b>bailu-ZZ</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/coldbrewtea">
        <img src="https://avatars.githubusercontent.com/u/6879314?v=4" width="80px" alt="coldbrewtea" />
        <br /><sub><b>coldbrewtea</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/excniesNIED">
        <img src="https://avatars.githubusercontent.com/u/61446002?v=4" width="80px" alt="excniesNIED" />
        <br /><sub><b>excniesNIED</b></sub>
      </a>
    </td>
  </tr>
  <tr>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/samz406">
        <img src="https://avatars.githubusercontent.com/u/3055810?v=4" width="80px" alt="samz406" />
        <br /><sub><b>samz406</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/imaben">
        <img src="https://avatars.githubusercontent.com/u/3390195?v=4" width="80px" alt="imaben" />
        <br /><sub><b>imaben</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/wangle201210">
        <img src="https://avatars.githubusercontent.com/u/19949348?v=4" width="80px" alt="wangle201210" />
        <br /><sub><b>wangle201210</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/CLukeLi">
        <img src="https://avatars.githubusercontent.com/u/252523101?v=4" width="80px" alt="CLukeLi" />
        <br /><sub><b>CLukeLi</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/myssl">
        <img src="https://avatars.githubusercontent.com/u/27838738?v=4" width="80px" alt="myssl" />
        <br /><sub><b>myssl</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/exgliuzhi">
        <img src="https://avatars.githubusercontent.com/u/6261701?v=4" width="80px" alt="exgliuzhi" />
        <br /><sub><b>exgliuzhi</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/hoorayman">
        <img src="https://avatars.githubusercontent.com/u/73151874?v=4" width="80px" alt="hoorayman" />
        <br /><sub><b>hoorayman</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/debin-ge">
        <img src="https://avatars.githubusercontent.com/u/21329997?v=4" width="80px" alt="debin-ge" />
        <br /><sub><b>debin-ge</b></sub>
      </a>
    </td>
  </tr>
  <tr>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/ocass-chen">
        <img src="https://avatars.githubusercontent.com/u/172055494?v=4" width="80px" alt="ocass-chen" />
        <br /><sub><b>ocass-chen</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/AnxForever">
        <img src="https://avatars.githubusercontent.com/u/130662349?v=4" width="80px" alt="AnxForever" />
        <br /><sub><b>AnxForever</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/yujiewanwan">
        <img src="https://avatars.githubusercontent.com/u/268286250?v=4" width="80px" alt="yujiewanwan" />
        <br /><sub><b>yujiewanwan</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/lxm">
        <img src="https://avatars.githubusercontent.com/u/1918195?v=4" width="80px" alt="lxm" />
        <br /><sub><b>lxm</b></sub>
      </a>
    </td>
    <td align="center" valign="top" width="12.5%">
      <a href="https://github.com/susunola">
        <img src="https://avatars.githubusercontent.com/u/38539169?v=4" width="80px" alt="susunola" />
        <br /><sub><b>susunola</b></sub>
      </a>
    </td>
  </tr>
</table>

<!-- readme: contributors -end -->

<p align="center">
  <a href="https://github.com/astaxie/TokenHub/graphs/contributors">Посмотреть всех участников</a>
  ·
  <a href="CONTRIBUTING.md">Принять участие в разработке</a>
</p>

## История звёзд (Star History)

<a href="https://www.star-history.com/?repos=astaxie%2Ftokenhub&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=astaxie/tokenhub&type=date&theme=dark&legend=top-left&sealed_token=hWH3kDnssTf49CCLxzq3NVqEp0WTL-HFhsdpQJJz1DUuZt0D-nu1jgXLnhCxrUrMYujv6IJJk12B1wCp5qiU2bU_J03ECSYvb3Y9Pv-gqX7RuwS4SehRrQ" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=astaxie/tokenhub&type=date&legend=top-left&sealed_token=hWH3kDnssTf49CCLxzq3NVqEp0WTL-HFhsdpQJJz1DUuZt0D-nu1jgXLnhCxrUrMYujv6IJJk12B1wCp5qiU2bU_J03ECSYvb3Y9Pv-gqX7RuwS4SehRrQ" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=astaxie/tokenhub&type=date&legend=top-left&sealed_token=hWH3kDnssTf49CCLxzq3NVqEp0WTL-HFhsdpQJJz1DUuZt0D-nu1jgXLnhCxrUrMYujv6IJJk12B1wCp5qiU2bU_J03ECSYvb3Y9Pv-gqX7RuwS4SehRrQ" />
 </picture>
</a>

## Лицензия (License)

TokenHub распространяется под лицензией [Apache License 2.0](LICENSE).
