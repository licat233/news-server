import requests
import pandas as pd
from bs4 import BeautifulSoup
import time
from tqdm import tqdm
import oss2
import requests
from oss2.credentials import EnvironmentVariableCredentialsProvider
import os
import uuid
import threading
import mysql.connector

startAt = time.time()

headers = {
    # 'user-agent': 'Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36 Edg/110.0.1587.63',
     'user-agent': 'Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)'
}
proxy = {
    "http": "http://127.0.0.1:7890",  # 代理服务器的IP地址和端口号
    "https": "http://127.0.0.1:7890"  # 代理服务器的IP地址和端口号
}

baseUrl = "https://read01.com/zh-tw"

ossBaseUrl = "http://news-storage.oss-cn-hongkong.aliyuncs.com"

# 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量 OSS_ACCESS_KEY_ID 和 OSS_ACCESS_KEY_SECRET 。
auth = oss2.ProviderAuth(EnvironmentVariableCredentialsProvider())
# yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
# 填写Bucket名称。
bucket = oss2.Bucket(auth, 'http://oss-cn-hongkong.aliyuncs.com', 'news-storage')

def createDatabaseConn():
    # 连接MySQL数据库
    cnx = mysql.connector.connect(user='root', password='', host='localhost', database='news')
    return cnx

def extract_uid(url):
    slash_index = url.rfind("/")  # 查找最后一个斜杠的索引位置
    dot_index = url.rfind(".")  # 查找最后一个点的索引位置
    code = url[slash_index + 1:dot_index]  # 提取斜杠和点之间的部分
    return code


def queryArticleExist(uid):
    # 执行查询语句
    query = f"SELECT * FROM article WHERE raw_uid = '{uid}'"
    cnx = createDatabaseConn()
    cursor = cnx.cursor()
    cursor.execute(query)
    result = cursor.fetchone()

    # 关闭游标和数据库连接
    cursor.close()
    cnx.close()

    # 检查查询结果是否存在数据
    if result:
        return True
    else:
        return False


def queryImageExist(source):
    # 执行查询语句
    query = f"SELECT * FROM images WHERE source = '{source}'"
    cnx = createDatabaseConn()
    cursor = cnx.cursor()
    cursor.execute(query)
    result = cursor.fetchone()

    # 关闭游标和数据库连接
    cursor.close()
    cnx.close()

    # 检查查询结果是否存在数据
    if result:
        return True, result[2]
    else:
        return False, ""


def insertArticle(raw_uid, title, category, publish_at, author, cover, text_content, html_content):
    insert_query = "INSERT INTO article (uuid, raw_uid, title, category, publish_at, author, cover, text_content, html_content, tags) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
    uuid = str(uuid.uuid4().hex)
    tags = category
    data = (uuid, raw_uid, title, category, publish_at, author, cover, text_content, html_content, tags)
    cnx = createDatabaseConn()
    cursor = cnx.cursor()
    cursor.execute(insert_query, data)
    cnx.commit()
    # 关闭游标和数据库连接
    cursor.close()
    cnx.close()

def insertImage(source, src):
    # 数据不存在，执行保存或其他操作
    insert_query = "INSERT INTO images (source, src) VALUES (%s, %s)"
    data = (source, src)
    cnx = createDatabaseConn()
    cursor = cnx.cursor()
    cursor.execute(insert_query, data)
    cnx.commit()
    # 关闭游标和数据库连接
    cursor.close()
    cnx.close()


def get_max_page_record(session, category_id):
    params = {
        'page': 9999,
    }
    response = session.get(f'{baseUrl}/{category_id}/', headers=headers, params=params, proxies=proxy)
    response.raise_for_status()
    html = response.text
    soup = BeautifulSoup(html, 'html.parser')
    activeLiE = soup.find('ul', {"class": "pagination"}).find(
        "li", {"class": "active"})
    pageNumberV = activeLiE.find('a').text
    items = soup.find_all('div', {"class": "list-item"})

    try:
        pageNum = int(pageNumberV)
        record = len(items) + (pageNum-1)*25
        return pageNum, record
    except ValueError:
        print("无法将字符串转换为整数")
    return 0, 0


def get_page(session, category_id, page):
    params = {
        'page': page,
    }

    response = session.get(f'{baseUrl}/{category_id}/', headers=headers, params=params, proxies=proxy)
    response.raise_for_status()

    html = response.text
    soup = BeautifulSoup(html, 'html.parser')
    title_urls = []
    activeLiE = soup.find('ul', {"class": "pagination"}).find(
        "li", {"class": "active"})
    pageNumber = activeLiE.find('a').text
    if pageNumber != str(page):
        return title_urls, True

    items = soup.find_all('div', {"class": "list-item"})

    for item in items:
        titleE = item.find('div', {"class": "title"})
        title = titleE.text.replace('\xa0', '')
        article_path = titleE.find('a')['href']
        url = baseUrl+article_path
        # 判断文章是否已经存在数据库中
        uid = extract_uid(url)
        exist = queryArticleExist(uid)
        if exist:
            # 若已经存在，则不爬取该文章
            title_urls.append(("", "", "", ""))
            continue
        coverE = item.find('div', {"class": "image"})
        cover = coverE.find('img').get('data-src', "")
        title_urls.append((title, url, cover, uid))

    return title_urls, False


def get_content(session, url):
    response = session.get(url, headers=headers)
    response.raise_for_status()

    html = response.text
    soup = BeautifulSoup(html, 'html.parser')

    dt, meta = soup.find('div', {'class': 'meta'}).text.strip().split('\xa0')
    # 开始识别并移除广告
    # 根据CSS类名进行筛选
    ads = soup.find_all(class_='axslot')
    for ad in ads:
        ad.decompose()  # 直接移除广告标签及其内容
    contentE = soup.find('div', {'class': 'content'})
    # 更换图片
    imgsE = contentE.find_all('img')
    for imgE in imgsE:
        img_url = imgE.get('data-src', None)
        if img_url is None:
            imgE.decompose()  # 移除
            continue
        if not img_url.startswith('http'):
            # 如果src属性值相对路径，则需要加上文章所在目录的路径
            img_url = os.path.join(os.path.dirname(url), img_url)
        retry_count=0
        new_img_url = upload_to_oss(img_url, "images")
        while new_img_url == "" and retry_count < 3:
            retry_count += 1
            new_img_url = upload_to_oss(img_url, "images")
        if new_img_url == "":
            return "","","",""
        imgE['src'] = new_img_url
        imgE['data-src'] = new_img_url
        # image_info = {
        #     "source": img_url,
        #     "src": new_img_url,
        # }
        # image_infos.append(image_info)

    content = contentE.text.strip().replace('\xa0', '')
    htmlContent = contentE.prettify()

    return dt, meta, content, htmlContent


def generate_random_filename(filename):
    # 获取文件的扩展名
    extension = os.path.splitext(filename)[1]

    # 生成唯一的标识符
    random_string = str(uuid.uuid4().hex)

    # 拼接新的文件名
    new_filename = random_string + extension

    return new_filename


def upload_to_oss(image_url, dirname):
    # 判断该图片是否已存在于oss中
    exist, src = queryImageExist(image_url)
    if exist:
        return src
    # requests.get返回的是一个可迭代对象（Iterable），此时Python SDK会通过Chunked Encoding方式上传。
    # 填写网络流地址。
    input = requests.get(image_url)
    # 填写Object完整路径。Object完整路径中不能包含Bucket名称。

    newName = generate_random_filename(image_url)
    try:
        bucket.put_object(f'{dirname}/{newName}', input)
        src = f"{ossBaseUrl}/{dirname}/{newName}"
        insertImage(image_url, src)
        return src
    except oss2.exceptions.ServerError as e:
        print(f'上传 {newName} 图片失败，来源：{image_url}')
    return ""


def category_handler(category, category_id):
    retries = 3  # 设置重试次数
    retry_delay = 3  # 设置重试间隔（秒）

    def aticle_handler(title, url, cover, uid, callback):
        if title == "" or url == "" or cover == "" :
            if callback:
                callback()
            return
        for attempt in range(retries):
            try:
                dt, meta, textContent, htmlContent = get_content(session, url)
                try_count = 0
                ossCover = upload_to_oss(cover, "covers")
                while ossCover == "" and try_count < 3:
                    ossCover = upload_to_oss(cover, "covers")
                if ossCover == "":
                    continue
                author = meta.replace('來源：', '')
                insertArticle(uid, title ,category, dt, author, ossCover, textContent, htmlContent)
                break
            except requests.exceptions.RequestException as e:
                print(f'获取 {category} 类别的文章失败：{title} (第 {attempt + 1} 次重试)')
                if attempt < retries - 1:
                    print(f'将在 {retry_delay} 秒后重试...')
                    time.sleep(retry_delay)
                else:
                    print(f'达到最大重试次数，放弃获取文章：{title}')
                    break
            except oss2.exceptions.ServerError as e:
                print(f'上传 {uid} 文章的图片失败：{title} (第 {attempt + 1} 次重试)')
                if attempt < retries - 1:
                    print(f'将在 {retry_delay} 秒后重试...')
                    time.sleep(retry_delay)
                else:
                    print(f'达到最大重试次数，放弃上传图片：{title}')
                    break
        if callback:
            callback()

    with requests.Session() as session:
        max_pagenum, max_item = get_max_page_record(session, category_id)
        # 创建进度条对象
        progress_bar = tqdm(total=max_item)
        def bar_add():
            progress_bar.update(1)
        for page in range(1, max_pagenum+1):
            title_urls, overflow = get_page(session, category_id, page)
            if overflow:
                break
            # for title, url, cover, uid in title_urls:
            #     aticle_handler(title, url, cover, uid, bar_add)
            # 使用多线程
            ts = []
            for title, url, cover, uid in title_urls:
                if title == "" or url == "" or cover == "" :
                    bar_add()
                    continue
                t = threading.Thread(target=aticle_handler, args=(title, url, cover, uid, bar_add), name=uid)
                ts.append(t)
            for t in ts:
                t.start()
            for t in ts:
                t.join()
        progress_bar.close()
    # df = pd.DataFrame(post_infos)
    # df.to_excel(f'{category}.xlsx', index=False)

categories = [('健康', 'c508'), ('美食', 'c710'), ('時尚', 'c942'), ('运势', 'c313'), ('娛樂', 'c475'), ('旅游', 'c808'), ('動漫', 'c921'), ('家居', 'c515'), ('科學', 'c610'), ('文化', 'c502'), ('歷史', 'c616')]
# categories = [('健康', 'c508')]

# ts=[]
for category, category_id in categories:
    category_handler(category, category_id)
#     t = threading.Thread(target=category_handler, args=(category, category_id), name=category_id)
#     ts.append(t)

# # 启动线程
# for t in ts:
#     t.start()

# # 等待线程执行完毕
# for t in ts:
#     t.join()

# print('All threads done!')

endAt = time.time()

subTime = endAt - startAt

print("耗时: ", subTime, "秒")