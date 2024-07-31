import json
import csv
import uuid
from faker import Faker
import random
import subprocess
import os
from datetime import datetime

# Инициализация Faker для генерации случайных данных
fake = Faker()

def generate_random_data():
    # Генерация случайных данных
    order_uid = str(uuid.uuid4())
    track_number = fake.uuid4()
    entry = random.choice(["WBIL", "WB", "WBILM"])
    delivery = {
        "name": fake.name(),
        "phone": fake.phone_number(),
        "zip": fake.zipcode(),
        "city": fake.city(),
        "address": fake.address(),
        "region": fake.state(),
        "email": fake.email()
    }
    payment = {
        "transaction": order_uid,
        "request_id": "",
        "currency": "USD",
        "provider": random.choice(["wbpay", "paypal"]),
        "amount": random.randint(1000, 2000),
        "payment_dt": int(datetime.now().timestamp()),
        "bank": random.choice(["alpha", "beta"]),
        "delivery_cost": random.randint(1000, 2000),
        "goods_total": random.randint(100, 500),
        "custom_fee": 0
    }
    items = [{
        "chrt_id": random.randint(1000000, 9999999),
        "track_number": track_number,
        "price": random.randint(100, 500),
        "rid": str(uuid.uuid4()),
        "name": fake.word(),
        "sale": random.randint(0, 50),
        "size": str(random.randint(0, 10)),
        "total_price": random.randint(100, 500),
        "nm_id": random.randint(1000000, 9999999),
        "brand": fake.company(),
        "status": random.choice([200, 201, 202])
    }]
    data = {
        "order_uid": order_uid,
        "track_number": track_number,
        "entry": entry,
        "delivery": delivery,
        "payment": payment,
        "items": items,
        "locale": "en",
        "internal_signature": "",
        "customer_id": "test",
        "delivery_service": "meest",
        "shardkey": "9",
        "sm_id": 99,
        "date_created": datetime.now().isoformat(),
        "oof_shard": "1"
    }
    return data, order_uid

def save_to_csv(order_uid):
    file_exists = os.path.isfile('order_uids.csv')

    with open('order_uids.csv', mode='a', newline='') as file:
        writer = csv.writer(file)

        # Записываем заголовки, если файл новый
        if not file_exists:
            writer.writerow(['order_uid'])

        # Записываем данные
        writer.writerow([order_uid])

def send_to_nats(data):
    # Запишите JSON данные в файл
    with open('model.json', 'w') as f:
        json.dump(data, f)

    # Выполните команду PowerShell для отправки данных в NATS
    command = "powershell.exe -Command \"Get-Content model.json | nats pub l0.test -s localhost:4223\""
    subprocess.run(command, shell=True)

def main():
    for _ in range(50):  
        data, order_uid = generate_random_data()
        save_to_csv(order_uid)
        send_to_nats(data)

if __name__ == "__main__":
    main()
