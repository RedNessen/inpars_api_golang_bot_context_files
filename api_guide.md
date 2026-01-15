Общая информация
При помощи API Вы можете легко и быстро интегрировать функционал нашего сервиса в свою систему. Информация отдается в двух форматах (json и xml), в зависимости от установленного заголовка (Accept) в запросе или при указании в GET параметре (_format). Для того чтобы начать работать с API, Вам необходимо зарегистрироваться, оплатить любой из тарифов, где указано API, скопировать токен (специальный код для доступа к API, найти его можно на странице профиля вашего аккуанта, вкладка API) и использовать его во всех запросах к сервису.

Запросы и ограничения
Внимание: частота запросов к нашему сервису ограничена (10 запросов в минуту), в тестовом режиме разрешено делать 1 запрос в секунду. Если Вы превысите это ограничение - сервер вернет в ответ HTTP-заголовок со статусом 429 Too Many Requests, немного подождав, Вы сможете возобновить работу.

С каждым ответом сервер отправляет HTTP-заголовки, а так же в самом теле ответа содержатся дополнительные информационные поля.

HTTP-заголовки

X-Rate-Limit-Limit максимальное количество запросов
X-Rate-Limit-Remaining оставшееся количество запросов
X-Rate-Limit-Reset количество секунд до полного восстановления лимита
Поле meta

Некоторые значения могут появляться только в определенных запросах.

limit количество элементов в запросе (равен значению параметра limit списка объявлений)
totalCount общее количество элементов
updateLimit доступное количество обновлений (1000 в месяц)
updateRemaining оставшееся количество обновлений на текущий момент
rateLimit максимальное количество запросов (аналог X-Rate-Limit-Limit)
rateRemaining оставшееся количество запросов (аналог X-Rate-Limit-Remaining)
rateReset количество секунд до полного восстановления лимита (аналог X-Rate-Limit-Reset)
Пример ответа с ошибкой:

HTTP/1.1 429 Too Many Requests
Date: Mon, 29 Jun 2020 07:50:54 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 0
X-Rate-Limit-Reset: 60
Content-Type: application/json; charset=UTF-8

{
	"name": "Too Many Requests",
	"message": "Rate limit exceeded.",
	"code": 0,
	"status": 429
}
Аутентификация
Для взаимодействия с API необходима аутентификация в системе одним из двух методов с использованием токена доступа (найти его можно на странице профиля вашего аккуанта).
В примерах ниже будем использовать токен aEcS9UfAagInparSiv23aoa_vPzxqWvm, так же Вы можете использовать его при тестировании. В разделе списка объявлений отображается срез неактуальных объектов (некоторые данные по полям скрыты или заменены на тестовые), ограничения по другим разделам отсутствуют.

HTTP Basic Auth
Для использования данного метода необходимо отправлять с запросами HTTP-заголовок Authorization: Basic с base64-кодированным значением username:. В качестве имени пользователя (username) используется токен доступа к API, символ : так же должен быть закодирован.

curl -X GET -H "Accept: application/json" -H "Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06" "https://inpars.ru/api/v2/estate"
Параметр запроса
Все запросы дополняются GET параметром access-token, значение которого равно токену для доступа к API.

curl -X GET -H "Accept: application/json" "https://inpars.ru/api/v2/estate?access-token=aEcS9UfAagInparSiv23aoa_vPzxqWvm"
Список объявлений
Возвращает список последних объявлений, отсортированных по дате изменения в порядке убывания, включая все оплаченные по тарифу регионы, с возможностью фильтрации по указанным ниже параметрам. Количество объявлений в одном запросе по умолчанию равно 500.

Примечание: каждый день в 00:00 GMT +03:00 происходит смещение объявлений (минус год от начала текущего дня), т.е. Вы можете получить список объявлений только за год.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/estate

Параметр	Тип	Описание
sortBy	string	Сортировка объектов. Возможные значения:
updated_desc - по убыванию даты изменения (задано по умолчанию)
updated_asc - по возрастанию даты изменения
created_desc - по убыванию даты создания
created_asc - по возрастанию даты создания
id_desc - по убыванию идентификатора объявления
id_asc - по возрастанию идентификатора объявления.
lastId	integer	Будут возвращены объекты с идентификатором большим, чем указанный. При использовании параметра, параметр sortBy доступен в двух значениях: id_desc (задано по умолчанию) и id_asc, остальные значения сортировки игнорируются.
timeStart	integer	Дата начала выборки в формате UNIX-time, используется поле updated (по умолчанию равно минус год от начала текущего дня).
timeEnd	integer	Дата конца выборки в формате UNIX-time, используется поле updated, не может быть меньше timeStart.
regionId	integer/string	Идентификатор региона. Можно запрашивать несколько регионов через запятую, например: regionId=1,2.
cityId	integer/string	Идентификатор города. Можно запрашивать несколько городов через запятую, например: cityId=1,2.
metroId	integer/string	Идентификатор метро. Можно запрашивать несколько метро через запятую, например: metroId=1,2.
typeAd	integer/string	Идентификатор типа недвижимости. Возможные значения:
1 - сдам
2 - продам
3 - сниму
4 - куплю.
Можно запрашивать несколько типов через запятую, например: typeAd=1,3.
sectionId	integer/string	Идентификатор раздела недвижимости. Можно запрашивать несколько разделов через запятую, например: sectionId=1,4
categoryId	integer/string	Идентификатор категории недвижимости. Можно запрашивать несколько категорий через запятую, например: categoryId=1,2.
sellerType
new	integer/string	Продавец. Возможные значения:
1 - собственник (задано по умолчанию)
2 - агент
3 - застройщик.
Можно запрашивать несколько через запятую, например: sellerType=1,2,3.
withAgent
deprecated	integer	0 - отображать только собственников (задано по умолчанию), 1 - отображать собственников, агентов, застройщиков, 2 - отображать только агентов, 3 - отображать только застройщиков.

Обратите внимание, данный параметр устарел и будет удален из будущих версий API!
withPhoto	integer	Отображать объявления 0 - только без фото, 1 - только с фото. Если параметр не указан, отображаются все объявления.
isNew	integer	Отображать 0 - вторичное жилье, 1 - новостройки. Если параметр не указан, отображаются все объявления.
costMin	integer	Цена от.
costMax	integer	Цена до.
floorMin	integer	Этаж от.
floorMax	integer	Этаж до.
sqMin	float	Площадь от (десятичный разделитель - точка или запятая).
sqMax	float	Площадь до (десятичный разделитель - точка или запятая).
sqLandMin	float	Площадь земли от (десятичный разделитель - точка или запятая).
sqLandMax	float	Площадь земли до (десятичный разделитель - точка или запятая).
sourceId	integer/string	Идентификатор источника. Возможные значения:
1 - avito.ru
2 - cian.ru
5 - youla.io
7 - sob.ru
9 - bazarpnz.ru
11 - move.ru
13 - realty.yandex.ru
19 - gipernn.ru
21 - orsk.ru
22 - domclick.ru
23 - doska.ykt.ru.
Можно запрашивать несколько источников через запятую, например: sourceId=1,2.
polygon	string	Поиск по полигону (области) координат, список из точек "широта" и "долгота". Коордианты передаются в виде lat1,lng1,lat2,lng2...
Если полигон с большим числом точек, рекомендуем передавать параметры в POST-запросе.

Пример, поиск объектов по адресу "Москва, Береговой проезд" &polygon=55.7585,37.50965,55.7585,37.50965,55.75842,37.50933,55.75828,37.50901,55.75813,37.50874,55.75794,37.50859,55.75776,37.50843,55.75758,37.50834,55.75739,37.50833,55.75721,37.50838,55.75706,37.50869,55.75696,37.50902,55.75684,37.50934,55.75679,37.50966,55.75671,37.50998,55.75657,37.51027,55.75647,37.51062,55.75641,37.51094,55.75637,37.51126,55.7564,37.51158,55.7565,37.51191,55.75668,37.51216,55.75685,37.51239,55.75702,37.51268,55.75719,37.51294,55.75737,37.51306,55.75755,37.51312,55.75771,37.51284,55.75785,37.51259,55.75799,37.51228,55.75816,37.51198,55.75831,37.51168,55.75839,37.51136,55.75846,37.51103,55.75848,37.51071,55.75851,37.51039,55.75854,37.51007,55.75855,37.50975

Обратите внимание, данный параметр не работает с тестовым токеном!
parseId	integer/string	Поиск по идентификатору объявления, используемого на источнике. Можно запрашивать несколько объектов через запятую (максимум 500), например: parseId=1234567890,123ab4567c.
Идентификатор может совпадать на разных источниках. Для получения точного результата используйте совместно с параметром sourceId.

Примечание: при использовании параметра учитывайте, что если в запросе присутствуют другие параметры, то они ограничивают видимость поиска или расширяют, например запрос с параметрами: parseId=...&sourceId=...&sellerType=1,2 отобразит собственников и агентов (указывайте sellerType=1,2,3 если хотите получить все данные, т.к. по умолчанию возвращаются только собственники).

Обратите внимание, данный параметр не работает с тестовым токеном!
byUrl	string	Поиск по ссылке на источник объявления.

Требования к значению:
должно быть закодировано (URL encoding)
должно соответствовать структуре URL
<схема>:[//<хост>[:<порт>]][/][?<параметры>][#<якорь>]
длина не должна превышать 1000 символов.
Например:
требуется найти объявление https://avito.ru/moskva/kvartiry/1-k_kvartira_11_m_914_et._1234567890
кодируем (в PHP можно использовать функцию urlencode) и подставляем в параметр byUrl=https%3A%2F%2Favito.ru%2Fmoskva%2Fkvartiry%2F1-k_kvartira_11_m_914_et._1234567890.

При использовании данного параметра, параметры sourceId и parseId будут проигнорированы.

Рекомендуем пользоваться поиском по parseId, т.к. он менее требователен.

Примечание: при использовании параметра учитывайте, что если в запросе присутствуют другие параметры, то они ограничивают видимость поиска или расширяют, т.е. запрос с параметрами: byUrl=...&sellerType=1,2 отобразит собственников и агентов (указывайте sellerType=1,2,3 если хотите получить все данные, т.к. по умолчанию возвращаются только собственники).

Обратите внимание, данный параметр не работает с тестовым токеном!
fields	string	Поля объявления, которые необходимо вернуть. Через запятую, без лишних символов. По умолчанию возвращаются все поля. Возможные значения: id, regionId, cityId, metroId, typeAd, sectionId, categoryId, title, address, floor, floors, sq, sqLand, sqLiving, sqKitchen, cost, text, images, lat, lng, name, phones, url, agent, source, sourceId, created, updated.
expandnew	string	Дополнительные поля, не включенные по умолчанию. Через запятую, без лишних символов. Возможные значения: region, city, type, section, category, metro, material, rentTime, isNew, rooms, history, phoneProtected, parseId, isApartments, landStatus, landCategory, rentTerms, house, sourceCreated, sourceUpdated, isAuction.
limit	integer	Ограничение на количество получаемых объектов (по умолчанию 500). Пример: limit=100, т.е. получить 100 объектов в одном запросе. Максимальное значение равно 1000, в тестовом режиме 50. Если в вашем регионе мало объектов, рекомендуем оставить значение по умолчанию, для максимальной скорости выборки.
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Обход коллекции объектов

Для обхода коллекции объектов используйте параметры timeStart, timeEnd и sortBy. Если в запросе объектов больше установленного лимита (поле totalCount больше поля limit), то используйте дату объекта с максимальным или минимальным значением поля updated (зависит от сортировки).

Предположим у Вас сортировка по возрастанию даты изменения sortBy=updated_asc, в таком случае нужно взять поле updated объекта с максимальным значением из списка (это последний объект) и подставить ее в следующий запрос в параметр timeStart (предварительно преобразовав время в формат UNIX-time), параметр timeEnd можно не указывать или указать значение текущей временной метки UNIX.

Предположим у Вас сортировка заданная по умолчанию, т.е. по убыванию даты изменения sortBy=updated_desc, в таком случае нужно взять поле updated объекта с минимальным значением из списка (это последний объект) и подставить ее в следующий запрос в параметр timeEnd (предварительно преобразовав время в формат UNIX-time), параметр timeStart указывать не нужно.

Рекомендуем использовать сортировку sortBy=updated_asc, так как она более привычна для обхода коллекции и временные интервалы в Вашей обработке будут последовательны.

Обратите внимание, что используя дату последнего полученного объекта и подставляя ее в новый запрос, Вы так же получите объекты с этой датой из предыдущего запроса, поэтому используйте фильтрацию на своей стороне.

Пример запроса на получение списка объявлений на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/estate', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTPHEADER => array(
        'Accept: application/json',
        'Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06',
    ),
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 07:12:40 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 106876174,
            "regionId": 77,
            "cityId": 1,
            "metroId": 30,
            "typeAd": 1,
            "sectionId": 6,
            "categoryId": 29,
            "title": "2-к квартира, 55 м², 9/14 эт.",
            "address": "Москва, Открытое ш., 24к11",
            "floor": 9,
            "floors": 14,
            "sq": 55,
            "sqLand": 0,
            "cost": 43000,
            "text": "Сдаю 2 комнатную квартиру 55 кв м., 9 этаж (14 эт.) панельного дома, комнаты изолированные (17 и 13 кв. м), кухня 9 кв. м. Вся мебель и бытовая техника есть. Окна выходят во двор, домофон, развитая инфраструктура. Собственник.",
            "images": [
                "https://22.img.avito.st/640x480/8845088222.jpg",
                "https://63.img.avito.st/640x480/8845089163.jpg",
                "https://31.img.avito.st/640x480/8845089631.jpg",
                "https://01.img.avito.st/640x480/8845090401.jpg",
                "https://56.img.avito.st/640x480/8845091056.jpg",
                "https://97.img.avito.st/640x480/8845091497.jpg"
            ],
            "lat": "55.823014",
            "lng": "37.757557",
            "name": "Оксана",
            "phones": [
                79000000000
            ],
            "url": "https://www.avito.ru/moskva/kvartiry/",
            "agent": 0,
            "source": "avito.ru",
            "sourceId": 1,
            "created": "2020-06-29T10:57:34+03:00",
            "updated": "2020-06-29T10:57:34+03:00"
        },
        {
            "id": 104873288,
            "regionId": 77,
            "cityId": 1,
            "metroId": 362,
            "typeAd": 1,
            "sectionId": 6,
            "categoryId": 32,
            "title": "Комната, 12.8 м²",
            "address": "Москва, Некрасовка, Некрасовская улица, 7",
            "floor": 16,
            "floors": 17,
            "sq": 53,
            "sqLand": 0,
            "cost": 15000,
            "text": "Комната очень уютная и теплая для одного человека очень хорошо.если для двоих то сумма будет на тысячу больше.рядом в шаговой доступности магазины аптеки сбербанк.и т.д",
            "images": [
                "https://cdn0.youla.io/files/images/orig/5e/da/5eda159cdaddaa6da4475084-1.jpg",
                "https://cdn0.youla.io/files/images/orig/5e/da/5eda15a195494856b062a598-1.jpg",
                "https://cdn0.youla.io/files/images/orig/5e/da/5eda15a5a38a1930b614e535-1.jpg",
                "https://cdn0.youla.io/files/images/orig/5e/da/5eda15a960c15645d65e1512-1.jpg"
            ],
            "lat": "55.681185",
            "lng": "37.924486",
            "name": "Ольга .",
            "phones": [
                79000000000
            ],
            "url": "//youla.io/moskva/nedvijimost/arenda-komnati/",
            "agent": 0,
            "source": "youla.io",
            "sourceId": 5,
            "created": "2020-06-05T12:53:24+03:00",
            "updated": "2020-06-29T10:57:34+03:00"
        },
        ...
    ],
    "meta": {
        "limit": 10,
        "totalCount": 4784081,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Объявление
Возвращает данные объявления по его идентификатору. В большинстве случаев, запрос использовать не обязательно, т.к. в списке объявлений отображаются все поля.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/estate/:id

Параметр	Тип	Описание
:id	integer	Идентификатор объявления (обязательный параметр).
fields	string	Поля объявления, которые необходимо вернуть. Через запятую, без лишних символов. По умолчанию возвращаются все поля. Возможные значения: id, regionId, cityId, metroId, typeAd, sectionId, categoryId, title, address, floor, floors, sq, sqLand, sqLiving, sqKitchen, cost, text, images, lat, lng, name, phones, url, agent, source, sourceId, created, updated.
expandnew	string	Дополнительные поля, не включенные по умолчанию. Через запятую, без лишних символов. Возможные значения: region, city, type, section, category, metro, material, rentTime, isNew, rooms, history, phoneProtected, parseId, isApartments, landStatus, landCategory, rentTerms, house, sourceCreated, sourceUpdated, isAuction.
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получения информации по объявлению с дополнительными полями на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/estate/106876174?expand=region,city,type,section,category,metro,material,rentTime,isNew,rooms,history,phoneProtected,parseId', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTPHEADER => array(
        'Accept: application/json',
        'Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06',
    ),
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 07:12:40 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": {
        "id": 106876174,
        "regionId": 77,
        "cityId": 1,
        "metroId": 30,
        "typeAd": 1,
        "sectionId": 6,
        "categoryId": 29,
        "title": "2-к квартира, 55 м², 9/14 эт.",
        "address": "Москва, Открытое ш., 24к11",
        "floor": 9,
        "floors": 14,
        "sq": 55,
        "sqLand": 0,
        "sqLiving": 32,
        "sqKitchen": 6,
        "cost": 43000,
        "text": "Сдаю 2 комнатную квартиру 55 кв м., 9 этаж (14 эт.) панельного дома, комнаты изолированные (17 и 13 кв. м), кухня 9 кв. м. Вся мебель и бытовая техника есть. Окна выходят во двор, домофон, развитая инфраструктура. Собственник.",
        "images": [
            "https://22.img.avito.st/640x480/8845088222.jpg",
            "https://63.img.avito.st/640x480/8845089163.jpg",
            "https://31.img.avito.st/640x480/8845089631.jpg",
            "https://01.img.avito.st/640x480/8845090401.jpg",
            "https://56.img.avito.st/640x480/8845091056.jpg",
            "https://97.img.avito.st/640x480/8845091497.jpg"
        ],
        "lat": "55.823014",
        "lng": "37.757557",
        "name": "Оксана",
        "phones": [
            79000000000
        ],
        "url": "https://www.avito.ru/moskva/kvartiry/",
        "agent": 1,
        "source": "avito.ru",
        "sourceId": 1,
        "created": "2020-06-29T10:57:34+03:00",
        "updated": "2020-06-29T10:57:34+03:00",
        "region": "Москва",
        "city": "Москва",
        "type": "Сдам",
        "section": "Жилая недвижимость",
        "category": "2-к квартира",
        "metro": "Бульвар Рокоссовского",
        "material": "панельный",
        "rentTime": 1,
        "isNew": false,
        "rooms": 2,
        "history": [
            {
                "date": "2020-06-29T10:57:34+03:00",
                "cost": 43000,
                "phones": [
                    79000000001
                ],
                "phoneProtected": true
            }
        ],
        "phoneProtected": true,
        "parseId": "1234567890",
        "isApartments": true,
        "rentTerms": {
            "commission": 40,
            "commissionType": 1,
            "deposit": 30000,
            "utilities": 1,
            "utilitiesMeters": 1,
            "utilitiesPrice": 3000
        }
    },
    "meta": {
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор объявления.
regionId	integer	Идентификатор региона.
cityId	integer	Идентификатор города.
metroId	integer	Идентификатор метро.
typeAd	integer	Идентификатор типа недвижимости. Возможные значения: 1 - сдам, 2 - продам, 3 - сниму, 4 - куплю.
sectionId	integer	Идентификатор раздела недвижимости.
categoryId	integer	Идентификатор категории недвижимости.
title	string	Заголовок объявления (может быть пустым).
address	string	Адрес объявления.
floor	integer	Этаж.
floors	integer	Этажность.
sq	float	Площадь (м2).
sqLand	float	Площадь участка (сот.).
sqLiving	float	Жилая площадь (м2).

Примечание: поле отсутствует, если значение не задано.
sqKitchen	float	Площадь кухни (м2).

Примечание: поле отсутствует, если значение не задано.
cost	integer	Стоимость.
text	string	Текст объявления.
images	array	Ссылки на изображения.
lat	string	Широта - точка координат.
lng	string	Долгота - точка координат.
name	string	Имя пользователя, разместившего объявление (если Собственник, то поле не заполняется).
phones	array/null	Телефоны пользователя, разместившего объявление.

Примечание: поле может быть пустым. Это значит, что собираются объявления "без звонков" и объявления, где по какой-то причине, номер не может быть собран.
url	string	Ссылка на источник объявления.
agentnew	integer	Продавец. Возможные значения: 0 - собственник, 1 - агент, 2 - застройщик.
source	string	Источник объявления. Возможные значения: avito.ru, cian.ru, youla.io, sob.ru, bazarpnz.ru, move.ru, realty.yandex.ru, gipernn.ru, orsk.ru, domclick.ru, doska.ykt.ru.
sourceId	integer	Идентификатор источника. Возможные значения:
1 - avito.ru
2 - cian.ru
5 - youla.io
7 - sob.ru
9 - bazarpnz.ru
11 - move.ru
13 - realty.yandex.ru
19 - gipernn.ru
21 - orsk.ru
22 - domclick.ru
23 - doska.ykt.ru
created	string	Дата добавления объявления (на сервер).
updated	string	Дата последнего изменения объявления.
region	string	Наименование региона.
city	string	Наименование города.
type	string	Наименование типа недвижимости.
section	string	Наименование раздела недвижимости.
category	string	Наименование категории недвижимости.
metro	string	Наименование метро.
material	string	Материал дома.
rentTime	integer	Срок аренды. Возможные значения: 0 - не указан, 1 - на длительный срок, 2 - посуточно.
isNew	boolean	Новостройка или вторичка. Возможные значения: true - новостройка, false - вторичка.
rooms	integer	Количество комнат.
history	array	История изменений. Возвращает массив в виде списка полей:
date - дата изменения
cost - стоимость
phones - телефоны (массив)
phoneProtected - телефон защищен
phoneProtected	boolean	Телефон защищен (подменный). Возможные значения: true - защищен, false - указан реальный номер.
parseId	string	Идентификатор объявления используемый на источнике.
isApartmentsnew	boolean	Является ли объект апартаментами. Возможные значения: true - апартаменты, иначе параметр отсутствует.
landStatusnew	integer	Статус земельного участка. Возможные значения:
1 - ИЖС (Индивидуальное жилищное строительство)
2 - СНТ/ДНП (Садоводческое некоммерческое товарищество, дачное некоммерческое партнёрство)
3 - ЛПХ (Личное подсобное хозяйство)
4 - КФХ (Фермерское хозяйство)
5 - Промназначения
landCategorynew	integer	Категория земельного участка. Возможные значения:
1 - Земли сельскохозяйственного назначения
2 - Земли населенных пунктов
3 - Другие (промышленности, лесного и водного фондов, земли запаса)
rentTermsnew	object	Условия аренды (см. поля ниже).
commission	integer	Комиссия, которую заплатит арендатор.
commissionType	integer	Тип комиссии. Возможные значения: 1 - процент со сделки, 2 - фиксированная сумма.
Используйте, чтобы определить какой формат у поля commission % или ₽.
deposit	integer	Залог.
utilities	integer	Оплата ЖКУ. Возможные значения: 1 - оплачивается арендатором, 2 - включена в платёж.
utilitiesMeters	integer	Оплата по счетчикам. Возможные значения: 1 - оплачивается арендатором, 2 - включена в платёж.
utilitiesPrice	integer	Стоимость ЖКУ.
housenew	object	Информация о доме
buildYear	integer	Год постройки
cargoLifts	integer	Количество грузовых лифтов
passengerLifts	integer	Количество пассажирских лифтов
sourceCreatednew	string	Дата публикации на источнике
sourceUpdatednew	string	Дата обновления на источнике
isAuctionnew	boolean	Объект выставлен на торги (аукцион). Возможные значения: true - аукцион, иначе параметр отсутствует.
Обновление объявлений
Позволяет обновить объекты или добавить пропущенные.

Объект будет добавлен/обновлен при следующих условиях:

если прошло 24 часа с момента последнего добавления/обновления объекта
если процесс обработки объекта не вызвал ошибку
Вы можете обновить до 1000 объектов в месяц (динамичное значение, зависит от нагрузки, смотрите значение updateLimit). Если в запросе присутствует ссылка на объект которого нет в базе данных, объект считается новым и не будет учитываться в лимите (значение updateRemaining не уменьшается).

Если Вы или другой пользователь в течении 24 часов отправили ссылку на добавление/обновление объекта, отправив повторный запрос с этой же ссылкой, будет возвращен результат с информацией о задаче. Информация будет доступна 24 часа, спустя 24 часа последующий запрос с этой же ссылкой создаст новую задачу на обновление.

Задача может выполняться в течении 24 часов, средний интервал добавления/обновления объекта составляет 5-10 минут.

На данный момент для добавления/обновления объектов доступны следующие источники: avito.ru, cian.ru, youla.io, gipernn.ru, domclick.ru, doska.ykt.ru (список будет расширяться), остальные источники будут проигнорированы.

Примечание: все возможные ограничения для запроса возможно будут изменены в будущем.

Внимание: на запрос не действуют общие ограничения на частоту запросов, что позволяет выполнять запросы к API параллельно с другими, разрешено делать 10 запросов в минуту.

Метод: POST

URL запроса: https://inpars.ru/api/v2/estate/task

Параметр	Тип	Описание
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Тело запроса: необходимо отправить данные в json формате с параметрами, указанными ниже

Параметр	Тип	Описание
urls	array	Список ссылок (максимум 10).
Каждый элемент списка должен соответствовать структуре URL
<схема>:[//<хост>[:<порт>]][/][?<параметры>][#<якорь>].
Пример запроса на добавление/обновление объектов на языке PHP:

<?php
				
// список ссылок
$urlList = [
    'urls' => [
        'https://avito.ru/moskva/kvartiry/1-k_kvartira_11_m_914_et._1234567890',
        'https://domclick.ru/card/sale__flat__1234567890'
    ]
];
				
$curl = curl_init();

curl_setopt_array($curl, [
    CURLOPT_URL => 'https://inpars.ru/api/v2/estate/task', // URL запроса, если требуется добавить ?_format=xml или ?_format=json
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTPHEADER => [
        'Accept: application/json',
        'Content-type: application/json',
        'Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06',
    ],
    CURLOPT_POST => true,
    CURLOPT_POSTFIELDS => json_encode($urlList), // кодируем в json
]);

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Tue, 30 May 2023 16:23:57 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 1,
            "url": "https://avito.ru/moskva/kvartiry/1-k_kvartira_11_m_914_et._1234567890",
            "isExist": false,
            "status": "created",
            "message": "URL добавлен в очередь.",
            "created": "2023-05-30T19:23:57+03:00",
            "updated": "2023-05-30T19:23:57+03:00"
        },
        {
            "id": 2,
            "url": "https://domclick.ru/card/sale__flat__1234567890",
            "isExist": true,
            "status": "inprogress",
            "message": "Объект в обработке.",
            "created": "2023-05-30T19:23:14+03:00",
            "updated": "2023-05-30T19:23:16+03:00"
        }
    ],
    "meta": {
        "updateLimit": 1000,
        "updateRemaining": 999,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор задачи.
url	string	Ссылка на объявление.
isExist	boolean	Информация о наличии объекта в базе данных. Возможные значения: true - существует, false - новый.
status	string	Статус задачи. Возможные значения:
error - ошибка
created - задача создана
inprogress - задача выполняется (может выполняться в течении 24 часов)
failed - задача не выполнена
completed - задача выполнена успешно.
message	string	Сообщение об ошибке или успехе выполнения задачи.
created	string	Дата создания задачи.
updated	string	Дата обновления задачи.
Дополнительные поля meta

updateLimit доступное количество обновлений (1000 в месяц)
updateRemaining оставшееся количество обновлений на текущий момент
Список регионов
Возвращает список регионов.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/region

Параметр	Тип	Описание
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получение списка регионов на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/region', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTPHEADER => array(
        'Accept: application/json',
        'Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06',
    ),
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 08:24:37 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 1,
            "title": "Республика Адыгея"
        },
        {
            "id": 2,
            "title": "Республика Башкортостан"
        },
        {
            "id": 3,
            "title": "Республика Бурятия"
        },
        ...
    ],
    "meta": {
        "totalCount": 85,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор региона.
title	string	Наименование региона.
Список городов
Возвращает список городов.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/city

Параметр	Тип	Описание
regionId	integer	Идентификатор региона. Если указан, то ответ будут отфильтрован по заданному региону.
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получение списка городов региона "Республика Адыгея" на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/city?regionId=1', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTPHEADER => array(
        'Accept: application/json',
        'Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06',
    ),
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 08:30:58 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 127,
            "title": "Майкоп",
            "regionId": 1
        },
        {
            "id": 128,
            "title": "Абадзехская",
            "regionId": 1
        },
        {
            "id": 129,
            "title": "Адыгейск",
            "regionId": 1
        },
        ...
    ],
    "meta": {
        "totalCount": 26,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор города.
title	string	Наименование города.
regionId	integer	Идентификатор региона.
Список метро
Возвращает список метро.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/metro

Параметр	Тип	Описание
regionId	integer	Идентификатор региона. Если указан, то ответ будут отфильтрован по заданному региону.
cityId	integer	Идентификатор города. Если указан, то ответ будут отфильтрован по заданному городу.
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получение списка метро региона "Москва" на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/metro?regionId=77', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTPHEADER => array(
        'Accept: application/json',
        'Authorization: Basic YUVjUzlVZkFhZ0lucGFyU2l2MjNhb2FfdlB6eHFXdm06',
    ),
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 08:35:17 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 1,
            "title": "Андроновка",
            "regionId": 77,
            "cityId": 1
        },
        {
            "id": 2,
            "title": "Авиамоторная",
            "regionId": 77,
            "cityId": 1
        },
        {
            "id": 3,
            "title": "Автозаводская",
            "regionId": 77,
            "cityId": 1
        },
        ...
    ],
    "meta": {
        "totalCount": 239,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор метро.
title	string	Наименование города.
regionId	integer	Идентификатор региона.
cityId	integer	Идентификатор города.
Список разделов
Возвращает список разделов недвижимости.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/estate/section

Параметр	Тип	Описание
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получение списка разделов недвижимости на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/estate/section?access-token=aEcS9UfAagInparSiv23aoa_vPzxqWvm', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
	CURLOPT_HTTPHEADER => array(
        'Accept: application/json'
	)
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 08:38:40 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 1,
            "typeId": 2,
            "title": "Жилая недвижимость"
        },
        {
            "id": 4,
            "typeId": 2,
            "title": "Коммерческая недвижимость"
        },
        {
            "id": 5,
            "typeId": 2,
            "title": "Загородная недвижимость"
        },
        {
            "id": 6,
            "typeId": 1,
            "title": "Жилая недвижимость"
        },
        {
            "id": 7,
            "typeId": 1,
            "title": "Коммерческая недвижимость"
        },
        {
            "id": 8,
            "typeId": 1,
            "title": "Загородная недвижимость"
        }
    ],
    "meta": {
        "totalCount": 6,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор раздела.
title	string	Наименование раздела
typeId	integer	Тип раздела недвижимости, не путать с typeAd (раздела списка объявлений). Возможные значения: 1 - аренда (сюда входят типы недвижимости: сдам и сниму), 2 - продажа (сюда входят типы недвижимости: продам и куплю).
Список категорий
Возвращает список категорий недвижимости.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/estate/category

Параметр	Тип	Описание
sectionId	integer	Идентификатор раздела. Если указан, то ответ будут отфильтрован по заданному разделу.
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получение списка категорий недвижимости на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/estate/category?access-token=aEcS9UfAagInparSiv23aoa_vPzxqWvm', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
	CURLOPT_HTTPHEADER => array(
        'Accept: application/json'
	)
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 08:40:17 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "id": 1,
            "title": "1-к квартира",
            "typeId": 2,
            "sectionId": 1
        },
        {
            "id": 2,
            "title": "2-к квартира",
            "typeId": 2,
            "sectionId": 1
        },
        {
            "id": 3,
            "title": "3-к квартира",
            "typeId": 2,
            "sectionId": 1
        },
        ...
    ],
    "meta": {
        "totalCount": 30,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
id	integer	Идентификатор категории.
title	string	Наименование категории
typeId	integer	Тип раздела недвижимости, не путать с typeAd (раздела списка объявлений). Возможные значения: 1 - аренда (сюда входят типы недвижимости: сдам и сниму), 2 - продажа (сюда входят типы недвижимости: продам и куплю).
sectionId	integer	Идентификатор раздела.
Подписка
Возвращает список активных подписок.

Методы: GET | POST

URL запроса: https://inpars.ru/api/v2/user/subscribe

Параметр	Тип	Описание
_format	string	Формат возвращаемых данных. Возможные значения: xml, json. По умолчанию необходимый формат будет определен на основе HTTP заголовка "Accept".
Пример запроса на получение списка активных подписок на языке PHP:

<?php
								
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'https://inpars.ru/api/v2/user/subscribe?access-token=aEcS9UfAagInparSiv23aoa_vPzxqWvm', // URL запроса
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_MAXREDIRS => 3,
    CURLOPT_TIMEOUT => 30,
	CURLOPT_HTTPHEADER => array(
        'Accept: application/json'
	)
));

$response = curl_exec($curl);
$error = curl_error($curl);

curl_close($curl);

if ($error) {
    echo 'cURL Error #:' . $error;
} else {
	echo $response;
}
Пример ответа в формате JSON:

HTTP/1.1 200 OK
Date: Mon, 29 Jun 2020 08:41:54 GMT
X-Rate-Limit-Limit: 10
X-Rate-Limit-Remaining: 9
X-Rate-Limit-Reset: 6
Content-Type: application/json; charset=UTF-8

{
    "data": [
        {
            "regionId": 55,
            "typeId": 1,
            "startTime": "2020-06-12T21:00:00+03:00",
            "endTime": "2020-07-13T21:00:00+03:00",
            "subscribe": "1 месяц",
            "api": false
        },
        {
            "regionId": 77,
            "typeId": 1,
            "startTime": "2020-03-02T10:10:00+03:00",
            "endTime": "2020-07-31T10:10:00+03:00",
            "subscribe": "API 1 месяц",
            "api": true
        }
    ],
    "meta": {
        "totalCount": 2,
        "rateLimit": 10,
        "rateRemaining": 9,
        "rateReset": 6
    }
}
Описание возвращаемых полей

Поле	Тип	Описание
regionId	integer	Идентификатор региона.
typeId	integer	Тип раздела недвижимости, не путать с typeAd (раздела списка объявлений). Возможные значения: 1 - аренда (сюда входят типы недвижимости: сдам и сниму), 2 - продажа (сюда входят типы недвижимости: продам и куплю).
startTime	string	Дата начала действия подписки.
endTime	string	Дата окончания действия подписки.
subscribe	string	Наименование подписки.
api	boolean	Тип подписки: сайт или API. Возможные значения: true - API, false - сайт.
Коды ошибок
Cписок кодов состояния HTTP, возвращаемых сервисом.

Код	Статус
200	OK. Все сработало именно так, как и ожидалось.
400	Неверный запрос. Может быть связано с разнообразными проблемами, такими как неправильные параметры действия, и т.д.
401	Аутентификация завершилась неудачно.
402	Для доступа к сервису необходима оплата.
403	Аутентифицированному пользователю не разрешен доступ к указанной точке входа API.
404	Запрошенные данные не существуют.
405	Метод не поддерживается. Сверьтесь со списком поддерживаемых HTTP-методов в заголовке Allow.
415	Не поддерживаемый тип данных. Запрашивается неправильный тип данных или номер версии.
429	Слишком много запросов. Запрос отклонен из-за превышения ограничения частоты запросов.
500	Внутренняя ошибка сервера.