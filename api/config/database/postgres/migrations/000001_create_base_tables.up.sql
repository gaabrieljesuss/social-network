CREATE TABLE usuarios
(
    id       int generated always as identity primary key,
    nome     varchar(50)  not null,
    nick     varchar(50)  not null unique,
    email    varchar(50)  not null unique,
    senha    varchar(100) not null,
    criadoEm timestamp default current_timestamp
);

CREATE TABLE seguidores
(
    usuario_id int not null,
    FOREIGN KEY (usuario_id)
    REFERENCES usuarios (id)
    ON DELETE CASCADE,
    seguidor_id int not null,
    FOREIGN KEY (seguidor_id)
    REFERENCES usuarios (id)
    ON DELETE CASCADE,
    primary key (usuario_id, seguidor_id)
);

CREATE TABLE publicacoes
(
    id       int generated always as identity primary key,
    titulo   varchar(50)  not null,
    conteudo varchar(300) not null,
    autor_id int          not null,
    FOREIGN KEY (autor_id)
    REFERENCES usuarios (id)
    ON DELETE CASCADE,
    curtidas int       default 0,
    criadaEm timestamp default current_timestamp
);

INSERT INTO usuarios (nome, nick, email, senha)
VALUES
    ('Usuário 1', 'usuario_1', 'usuario1@gmail.com', '$2a$10$Pf9VaZrjfAsF14i5Zmzcbe6WbLLPnFh0U7zQAc1V6r0bt/XVw8Ml.'),
    ('Usuário 2', 'usuario_2', 'usuario2@gmail.com', '$2a$10$Pf9VaZrjfAsF14i5Zmzcbe6WbLLPnFh0U7zQAc1V6r0bt/XVw8Ml.'),
    ('Usuário 3', 'usuario_3', 'usuario3@gmail.com', '$2a$10$Pf9VaZrjfAsF14i5Zmzcbe6WbLLPnFh0U7zQAc1V6r0bt/XVw8Ml.');

INSERT INTO seguidores(usuario_id, seguidor_id)
VALUES
    (1, 2),
    (3, 1),
    (1, 3);

INSERT INTO publicacoes(titulo, conteudo, autor_id)
VALUES
    ('Publicação do Usuário 1', 'Essa é a publicação do usuário 1! Oba!', 1),
    ('Publicação do Usuário 2', 'Essa é a publicação do usuário 2! Oba!', 2),
    ('Publicação do Usuário 3', 'Essa é a publicação do usuário 3! Oba!', 3);