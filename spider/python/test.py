import threading
import time
import uuid

# 定义需要多线程执行的函数
def print_num(num):
    time.sleep(num)
    print(threading.current_thread().name)
    print(str(uuid.uuid4().hex))

# 创建两个线程
t1 = threading.Thread(target=print_num, args=(1,), name='Thread1')
t2 = threading.Thread(target=print_num, args=(1,), name='Thread2')

# 启动线程
t1.start()
t2.start()

# 等待线程执行完毕
t1.join()
print("a")
t2.join()
print("b")

print('All threads done!')
