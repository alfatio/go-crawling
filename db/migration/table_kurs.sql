CREATE TABLE IF NOT EXISTS kurs (
    id int auto_increment primary key,
    symbol varchar(3) not null,
    e_rate_beli decimal(10,2) not null,
    e_rate_jual decimal(10,2) not null,
    tt_counter_beli decimal(10,2) not null,
    tt_counter_jual decimal(10,2) not null,
    bank_notes_beli decimal(10,2) not null,
    bank_notes_jual decimal(10,2) not null,
    kurs_date date not null
);

CREATE INDEX idx_kurs_date_symbol ON kurs (kurs_date, symbol);

CREATE UNIQUE INDEX unique_kurs ON kurs (symbol, kurs_date);