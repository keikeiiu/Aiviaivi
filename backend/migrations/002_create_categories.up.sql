CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO categories (name, slug) VALUES
    ('Animation', 'animation'),
    ('Comics', 'comics'),
    ('Gaming', 'gaming'),
    ('Music', 'music'),
    ('Dance', 'dance'),
    ('Knowledge', 'knowledge'),
    ('Tech', 'tech'),
    ('Sports', 'sports'),
    ('Lifestyle', 'lifestyle'),
    ('Movies', 'movies'),
    ('TV Shows', 'tv-shows'),
    ('Documentary', 'documentary');
