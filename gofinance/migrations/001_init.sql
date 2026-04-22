-- 001_init.sql

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('income', 'expense'))
);

CREATE TABLE IF NOT EXISTS transactions (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id),
    amount      NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    comment     TEXT NOT NULL DEFAULT '',
    date        DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date    ON transactions(date);

-- Стартовые категории
INSERT INTO categories (name, type) VALUES
    ('Зарплата',       'income'),
    ('Фриланс',        'income'),
    ('Подработка',     'income'),
    ('Прочий доход',   'income'),
    ('Продукты',       'expense'),
    ('Кафе и рестораны','expense'),
    ('Транспорт',      'expense'),
    ('Жильё и ЖКХ',   'expense'),
    ('Здоровье',       'expense'),
    ('Одежда',         'expense'),
    ('Развлечения',    'expense'),
    ('Связь',          'expense'),
    ('Образование',    'expense'),
    ('Прочие расходы', 'expense')
ON CONFLICT DO NOTHING;
