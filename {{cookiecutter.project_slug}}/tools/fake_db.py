"""
pip install psycopg2-binary faker pyyaml

python fake_db --verbose --bulk <table_name> <num_records>
"""

import argparse
import random

import psycopg2
import psycopg2.extras
import yaml
from faker import Faker

fake = Faker()

custom_fake_registry = []


def _get_fake_fn(field):
    if field['kind'] == "choice":
        def fn(column_name, data_type, character_maximum_length):
            if column_name == field["name"]:
                return random.choice(field['values']), True
            return None, False

        return fn

    raise Exception("Invalid kind")


def parse_custom_fake(fake_configs):
    for field in fake_configs['fields']:
        custom_fake_registry.append(_get_fake_fn(field))


def generate_fake_value(column_name, data_type, character_maximum_length):
    for fn in custom_fake_registry:
        v, fake_success = fn(column_name, data_type, character_maximum_length)
        if fake_success:
            return v

    if data_type in ['integer', 'bigint']:
        return fake.random_int(min=1, max=1000)
    elif data_type in ['varchar', 'character varying']:
        w = " ".join(fake.words(nb=int(character_maximum_length / 4)))
        return w[:character_maximum_length - 1]
    elif data_type == 'text':
        return fake.text()
    elif data_type in ['timestamp with time zone', 'timestamp without time zone']:
        return fake.date_time_this_decade()
    elif data_type == 'date':
        return fake.date()
    elif data_type == 'boolean':
        return fake.boolean()
    elif data_type == 'numeric':
        return fake.pydecimal(left_digits=5, right_digits=2, positive=True)
    elif data_type == 'uuid':
        return fake.uuid4()
    elif data_type == 'jsonb':
        return "{}"
    # Add more data type handlers as needed
    else:
        return fake.word()


def fake_data(conn, table_name, num_records, use_bulk=False, enable_debug=False):
    cursor = conn.cursor()
    cursor.execute(f"""
        SELECT column_name, data_type, character_maximum_length
        FROM information_schema.columns 
        WHERE table_name = '{table_name}';
    """)
    columns = cursor.fetchall()

    bulk_values = []
    current_index = 0
    for _ in range(num_records):
        values = [generate_fake_value(*c) for c in columns]
        placeholders = ', '.join(['%s'] * len(values))

        column_names_list = [column for column, _, __ in columns]
        column_names_str = ', '.join(column_names_list)

        current_index += 1
        if enable_debug:
            print(" | ".join([f"{c}:{v}" for c, v in zip(column_names_list, values)]))
            print("-" * 120)

        if use_bulk:
            bulk_values.append(tuple(values))
            if len(bulk_values) >= 1000:
                insert_query = f"INSERT INTO {table_name} ({column_names_str}) VALUES %s"
                psycopg2.extras.execute_values(cursor, insert_query, bulk_values)
                bulk_values = []
                conn.commit()
                print(f"finished: {current_index}")
        else:
            cursor.execute(
                f"INSERT INTO {table_name} ({column_names_str}) VALUES ({placeholders})",
                values
            )
            if current_index // 100 == 0:
                print(f"finished: {current_index}")

    if len(bulk_values) > 0:
        insert_query = f"INSERT INTO {table_name} ({column_names_str}) VALUES %s"
        psycopg2.extras.execute_values(cursor, insert_query, bulk_values)

    conn.commit()
    cursor.close()
    conn.close()
    print(f"Inserted {num_records} records into table {table_name}.")


def main():
    parser = argparse.ArgumentParser(description='Fake database.')
    parser.add_argument('table', type=str, help='Table name to fake')
    parser.add_argument('num', type=int, help='Num rows to fake')
    parser.add_argument('-v', '--verbose', action=argparse.BooleanOptionalAction, default=False, help='Enable debug')
    parser.add_argument('-b', '--bulk', action=argparse.BooleanOptionalAction, default=False, help='Bulk insert')

    # Parse the arguments
    args = parser.parse_args()

    with open("fake.yaml", 'r') as file:
        configs = yaml.safe_load(file)

    conn = psycopg2.connect(
        host=configs['db']['host'],
        port=configs['db']['port'],
        user=configs['db']['user'],
        password=configs['db']['pass'],
        dbname=configs['db']['name'],
    )

    parse_custom_fake(configs['fake'])
    fake_data(conn, args.table, args.num, args.bulk, args.verbose)


if __name__ == '__main__':
    main()
