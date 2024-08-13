"""
pip install psycopg2-binary faker

python fake_db <table_name> <num_records>
"""

from sys import argv

import psycopg2
from faker import Faker

conn = psycopg2.connect(
    dbname="{{ cookiecutter.project_slug }}",
    user="dev",
    password="dev",
    host="devhost",
    port="5432"
)

fake = Faker()


def generate_fake_value(column_name, data_type, character_maximum_length):
    if data_type in ['integer', 'bigint']:
        return fake.random_int(min=1, max=1000)
    elif data_type in ['varchar', 'character varying']:
        w = " ".join(fake.words(nb=int(character_maximum_length / 4)))
        return w[:character_maximum_length-1]
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


def fake_data(table_name, num_records):
    cursor = conn.cursor()
    cursor.execute(f"""
        SELECT column_name, data_type, character_maximum_length
        FROM information_schema.columns 
        WHERE table_name = '{table_name}';
    """)
    columns = cursor.fetchall()

    for _ in range(num_records):
        values = [generate_fake_value(*c) for c in columns]
        placeholders = ', '.join(['%s'] * len(values))
        column_names = ', '.join([column for column, _, __ in columns])

        print(values)
        cursor.execute(
            f"INSERT INTO {table_name} ({column_names}) VALUES ({placeholders})",
            values
        )

    conn.commit()
    cursor.close()
    conn.close()
    print(f"Inserted {num_records} records into table {table_name}.")


if __name__ == '__main__':
    table = argv[1]
    rows = int(argv[2])

    fake_data(table, rows)
